const net = require('net');
const tls = require('tls');
const fs = require('fs');

function connect(opts) {
  return new Promise((res, rej) => {
    let key, cert, ca
    try {
      key = fs.readFileSync(opts.key);
      cert = fs.readFileSync(opts.cert);
      ca = fs.readFileSync(opts.ca);
    } catch (err) {
      rej(err)
    }

    const tlsConn = tls.connect({
      key: key,
      cert: cert,
      ca: ca,
      host: opts.host,
      port: opts.port,
      servername: 'server.example.com'
    },
    () => {
      if (tlsConn.authorized) {
        return res(tlsConn)
      } else {
        const err = new Error('unauthorized');
        return rej(err)
      }
    })

      tlsConn.on('error', err => {
        return rej(err)
      })

      tlsConn.on('secureConnect', () => {
        return res(tlsConn)
      })
  })
}

module.exports = {
  connect
}