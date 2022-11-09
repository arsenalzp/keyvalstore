/**
 * outgoing message
 * 
 * cmd(3) - GET | SET | DEL | IMP | EXP
 * key(256)
 * value(512)
 * 
 * incoming message
 * 
 * error(1) - 79 | 78
 * data(512)
 */

const RESPBUF_SIZE = 512; // incoming message
const MESSAGE_SIZE = 771; // outgoing message

const OK = 79; // 'OK' code for server response
const CLOSED = 'readOnly' || 'writeOnly';
const GET_CMD = 'get';
const SET_CMD = 'set';
const DEL_CMD = 'del';
const IMP_CMD = 'imp';
const EXP_CMD = 'exp';

const { Readable } = require('stream');
const errors = require('./errors');
const clnt = require('./network');
const validateData = require('./utils');

class ReadableBuffer extends Readable {
  constructor(buffer) {
    super(buffer);
    this.buffer = buffer;
    this.counter = 0;
  }

  _read(size) {
    if (this.counter > this.buffer.length) {
      return null
    }

    if ((this.counter + size) >= this.buffer.length) {
      const slice = this.buffer.slice(this.counter, this.buffer.length);
      this.counter += size;
      this.push(Buffer.from(slice));
    } 

    const slice = this.buffer.slice(this.counter, (this.counter + size));
    this.counter += size;
    this.push(Buffer.from(slice));
  }
}

class KeyvalClnt {

  constructor(tlsConf) {
    this.ca = tlsConf.ca;
    this.key = tlsConf.key;
    this.cert = tlsConf.cert;
    this.host = tlsConf.host; 
    this.port = tlsConf.port;
    this.conn = null;
    this.mux = false;
  }

  async connect() {
    try {
      return this.conn = await clnt.connect({
        cert: this.cert, 
        key: this.key, 
        ca: this.ca,
        host: this.host,
        port: this.port,
      })
    } catch(err) {
      throw (err)
    }
    
  }

  lock() {
    this.mux = true;
  }

  unlock() {
    this.mux = false
  }

  isLocked() {
    return this.mux
  }
  async set(k, v) {
    try {
      while(this.isLocked()) {}
      this.lock()

      return await this._set(k, v)
    } catch(err) {
      throw err
    } finally {
      this.unlock()
    }
  }

  async get(k) {
    try {
      while(this.isLocked()) {}
      this.lock()

      return await this._get(k)
    } catch(err) {
      throw err
    } finally {
      this.unlock()
    }
  }

  async del(k) {
    try {
      while(this.isLocked()) {}
      this.lock()

      return await this._del(k)
    } catch(err) {
      throw err
    } finally {
      this.unlock()
    }
  }

  async imp(data) {
    try {
      while(this.isLocked()) {}
      this.lock()

      return await this._imp(data)
    } catch(err) {
      throw err
    } finally {
      this.unlock()
    }
  }

  async exp() {
    try {
      while(this.isLocked()) {}
      this.lock()

      return await this._exp()
    } catch(err) {
      throw err
    } finally {
      this.unlock()
    }
  }

  /**
   * Set the VALUE for the given KEY
   * 
   * @param {String} k a key 
   * @param {String} v a value 
   * @returns {Promise}
   */
  _set(k,v) {
    return new Promise((res, rej) => {

      try {
        validateData(k, v)
      } catch(err) {
        return rej(err)
      }

      const conn = this.conn;
      if (conn.readyState == CLOSED) {
        return rej(err)
      }

      conn.setKeepAlive(true)

      // create a buffer for outcoming messages
      const buf = Buffer.alloc(MESSAGE_SIZE);
      buf.write(SET_CMD);
      buf.write(k, 3);
      buf.write(v, 259);

      conn.write(buf);
      conn.setMaxListeners(2048);

      // create zero-length array for incoming data
      const respBuf = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();

        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == 0 || respBuf.length == RESPBUF_SIZE) {
            if (respBuf[0] != OK) {
              const respErr = respBuf.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res()
          }
          respBuf.push(byte)
        }
      });
    })
  }
  
  /**
   * Get the VALUE of the given KEY
   *
   * @param {String} k a key 
   * @returns {Promise} resolve Buffer
   */
  _get(k) {
    return new Promise((res, rej) => {

      try {
        validateData(k)
      } catch(err) {
        return rej(err)
      }

      const conn = this.conn;
      if (conn.readyState == CLOSED) {
        return rej()
      }

      conn.setKeepAlive(true);

      const buf = Buffer.alloc(MESSAGE_SIZE);
      buf.write(GET_CMD);
      buf.write(k, 3);

      conn.write(buf);
      conn.setMaxListeners(2048);
      
      // create zero-length array for incoming data
      const respBuf = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });
      
      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == 0 || respBuf.length == RESPBUF_SIZE) {
            if (respBuf[0] != OK) {
              const respErr = respBuf.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res(Buffer.from(respBuf.slice(1,)))
          }
          respBuf.push(byte)
        }
      });
    })
  }

  /**
   * Delete the give key
   *
   * @param {String} k a key
   * @returns {Promise}
   */
  _del(k) {
    return new Promise((res, rej) => {

      try {
        validateData(k)
      } catch(err) {
        return rej(err)
      }

      const conn = this.conn;
      if (conn.readyState == CLOSED) {
        return rej(err)
      }

      conn.setKeepAlive(true);

      const buf = Buffer.alloc(MESSAGE_SIZE);
      buf.write(DEL_CMD);
      buf.write(k, 3);

      conn.write(buf);
      conn.setMaxListeners(2048);

      // create zero-length array for incoming data
      const respBuf = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == 0 || respBuf.length == RESPBUF_SIZE) {
            if (respBuf[0] != OK) {
              const respErr = respBuf.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res()
          }
          respBuf.push(byte)
        }
      });
    })
  }

  /**
   * Import JSON data
   * 
   * @param {String|Buffer} data value 
   * @returns {Promise}
   */
  _imp(data) {
    return new Promise((res, rej) => {
      
      if (typeof data == 'object') {
        try {
          validateData(data);
        } catch(err) {
          return rej(err)
        } 
      } else if (typeof data == 'string') {
        const buf = Buffer.from(data);
        try {
          validateData(buf);
        } catch(err) {
          return rej(err)
        }
      } else {
        return rej(new Error(errors.UnknownTypeErr))
      }

      const conn = this.conn;
      if (conn.readyState == CLOSED) {
        return rej()
      }

      conn.setKeepAlive(true);

      // allocate the new buffer for command message
      const cmdBuf = Buffer.alloc(3);
      cmdBuf.write(IMP_CMD);
      
      // allocate the new buffer from data message and add a delimiter code ('\x00')
      const dataBuf = Buffer.from(data+'\x00', 'utf8');
      
      // if size of command and data buffers is less than MESSAGE_SIZE
      // create new buffer and fill necessary size up by zero
      if ((dataBuf.length + cmdBuf.length) < MESSAGE_SIZE) {
        const diffBuf = Buffer.alloc(MESSAGE_SIZE - (dataBuf.length + cmdBuf.length));
        const growBuf = Buffer.concat([cmdBuf, dataBuf, diffBuf]);
        const readable = new ReadableBuffer(growBuf);
        readable.pipe(conn)
      } else {
        const growBuf = Buffer.concat([cmdBuf, dataBuf]);
        const readable = new ReadableBuffer(growBuf);
        readable.pipe(conn);
      }

      conn.setMaxListeners(2048);

      // create zero-length array for incoming data
      const respBuf = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == 0 || respBuf.length == RESPBUF_SIZE) {
            if (respBuf[0] != OK) {
              const respErr = respBuf.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res()
          }
          respBuf.push(byte)
        }
      });
    })
  }

  /**
  * Export JSON data
  * 
  * @returns {Promise}
  */
  _exp() {
    return new Promise((res, rej) => {

      let data = Buffer.alloc(0);
      const conn = this.conn;
      if (conn.readyState == CLOSED) {
        return rej(err)
      }

      conn.setKeepAlive(true);

      const buf = Buffer.alloc(MESSAGE_SIZE);
      buf.write(EXP_CMD);

      conn.write(buf);
      conn.setMaxListeners(2048);
      
      // create zero-length array for incoming data
      const respBuf = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == 0) {
            if (respBuf[0] != OK) {
              const respErr = respBuf.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res(Buffer.from(respBuf.slice(1,)))
          }
          respBuf.push(byte)
        }
      });
    })
  }
}

Object.assign(KeyvalClnt, errors);
module.exports = KeyvalClnt;
