const axios = require('axios');

const baseURL = 'http://localhost:8080';
const apiKey = 'your-api-key';
const app = 'mytool';
let licenseKey = '';

async function createKey() {
  const res = await axios.post(`${baseURL}/create-key`, {
    app,
    duration: '30d'
  }, {
    headers: { 'X-Api-Key': apiKey }
  });
  console.log('✅ Create Key:', res.data);
  licenseKey = res.data.key;
}

async function login(hwid) {
  const res = await axios.post(`${baseURL}/login`, {
    app,
    key: licenseKey,
    hwid
  });
  console.log('✅ Login:', res.data);
}

async function deleteKey() {
  const res = await axios.post(`${baseURL}/delete-key`, {
    app,
    key: licenseKey
  }, {
    headers: { 'X-Api-Key': apiKey }
  });
  console.log('✅ Delete Key:', res.data);
}

(async () => {
  await createKey();
  await login('hwid-123');
  await login('hwid-123'); // second login
  await login('wrong-hwid'); // invalid hwid
  await deleteKey();
})();
