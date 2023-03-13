/**
 * outgoing message
 * 
 * cmd(3) - GET | SET | DEL | IMP | EXP
 * key(256)
 * value(511)
 * 
 * incoming message
 * 
 * error(1) - 79 | 78
 * data(512)
 */

const MESSAGE_SIZE = 771; // outgoing message
const EOT = '\u0004'; // End-Of-Trasmission character
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
      const dataBuffer = Buffer.alloc(MESSAGE_SIZE);
      dataBuffer.write(SET_CMD);
      dataBuffer.write(k, 3);
      dataBuffer.write(v, 259);
      dataBuffer.write(EOT, 770)

      conn.write(dataBuffer);
      conn.setMaxListeners(2048);

      // create zero-length array for incoming data
      const respBuffer = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();

        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == EOT.charCodeAt()) {
            if (respBuffer[0] == errors.ServerResponseError) {
              const respErr = respBuffer.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }
            
            return res()
          }
          respBuffer.push(byte)
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

      const dataBuffer = Buffer.alloc(MESSAGE_SIZE);
      dataBuffer.write(GET_CMD);
      dataBuffer.write(k, 3);
      dataBuffer.write(EOT, 770);

      conn.write(dataBuffer);
      conn.setMaxListeners(2048);
      
      // create zero-length array for incoming data
      const respBuffer = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });
      
      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == EOT.charCodeAt()) {
            if (respBuffer[0] == errors.ServerResponseError) {
              const respErr = respBuf.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res(Buffer.from(respBuffer.slice(1,)))
          }

          respBuffer.push(byte)
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

      const dataBuffer = Buffer.alloc(MESSAGE_SIZE);
      dataBuffer.write(DEL_CMD);
      dataBuffer.write(k, 3);
      dataBuffer.write(EOT, 770);

      conn.write(dataBuffer);
      conn.setMaxListeners(2048);

      // create zero-length array for incoming data
      const respBuffer = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == EOT.charCodeAt()) {
            if (respBuffer[0] == errors.ServerResponseError) {
              const respErr = respBuffer.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res()
          }
          respBuffer.push(byte)
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
      const dataSize = IMP_CMD.length + data.length + EOT.length;
      const dataBuffer = Buffer.alloc(dataSize,'utf8');
      dataBuffer.write(IMP_CMD, 0);
      dataBuffer.write(data + EOT, 3);
      const readable = new ReadableBuffer(dataBuffer);
      readable.pipe(conn);

      conn.setMaxListeners(2048);

      // create zero-length array for incoming data
      const respBuffer = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == EOT.charCodeAt()) {
            if (respBuffer[0] == errors.ServerResponseError) {
              const respErr = respBuffer.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res()
          }
          respBuffer.push(byte)
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

      const conn = this.conn;
      if (conn.readyState == CLOSED) {
        return rej(err)
      }

      conn.setKeepAlive(true);

      const dataBuffer = Buffer.alloc(MESSAGE_SIZE);
      dataBuffer.write(EXP_CMD);
      dataBuffer.write(EOT, 3);
      
      conn.write(dataBuffer);
  
      conn.setMaxListeners(2048);
      
      // create zero-length array for incoming data
      const respBuffer = [];

      conn.on('close', () => {
        return rej('closed')
      });

      conn.on('error', err => {
        conn.end();
        
        return rej(err)
      });

      conn.on('data', (buf) => {
        for (let byte of buf) {
          if (byte == EOT.charCodeAt()) {
            if (respBuffer[0] == errors.ServerResponseError) {
              const respErr = respBuffer.slice(1,);
              const err = new Error(respErr.toString());
    
              return rej(err)
            }

            return res(Buffer.from(respBuffer.slice(1,)))
          }
          respBuffer.push(byte)
        }
      });
    })
  }
}

Object.assign(KeyvalClnt, errors);
module.exports = KeyvalClnt;
