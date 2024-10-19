
const errors = require('../errors');

const KEY_SIZE = 256;
const VALUE_SIZE = 512;

function validateData(arg1, arg2) {
  switch (typeof arg1) {
    case 'object':
      try {
        JSON.parse(arg1.toString())
        break
      } catch(err) {
        throw err
      }
    case 'string':
      if (arg1.length > KEY_SIZE) {
        throw new Error(errors.KeySizeErr)
      }

      if (arg2 !== undefined && arg2.length > VALUE_SIZE) {
        throw new Error(errors.ValueSizeErr)
      }

      break;
  }
}

module.exports = validateData;