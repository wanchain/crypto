var ethUtil = require('ethereumjs-util');
var crypto = require('crypto');

function getPubKeyStrByPrivateKey(privateKey){
  if(!privateKey.startsWith('0x')){
    privateKey = '0x' + privateKey;
  }
  return ethUtil.bufferToHex(ethUtil.privateToPublic(privateKey), 'hex');
}

//generate random private key
function _generatePrivateKey(){
  var randomBuf = crypto.randomBytes(32);
  if (ethUtil.secp256k1.privateKeyVerify(randomBuf)){
    return randomBuf;
  } else {
    return _generatePrivateKey();
  }
}

function generateA1(RPrivateKeyDBytes, pubKeyA,  pubKeyB){
  A1 = ethUtil.secp256k1.publicKeyTweakMul(pubKeyA, RPrivateKeyDBytes, false);
  A1Bytes = ethUtil.sha3(A1);
  A1 = ethUtil.secp256k1.publicKeyTweakAdd(pubKeyB, A1Bytes, false);
  return A1;
}

function generateOTAPublicKey(pubKeyA, pubKeyB){
  RPrivateKey = _generatePrivateKey();
  A1 = generateA1(RPrivateKey, pubKeyA, pubKeyB);
  return {otaPubA1: A1, otaPubS1: ethUtil.secp256k1.publicKeyCreate(RPrivateKey)};
}

function computeOTAPrivateKeys(privatekey_a, privatekey_b, otaPubA1, otaPubS1){
  var pub = ethUtil.secp256k1.publicKeyTweakMul(ethUtil.secp256k1.publicKeyConvert(otaPubS1), privatekey_b, false);
  k = ethUtil.sha3(pub);
  k = ethUtil.secp256k1.privateKeyTweakAdd(k, privatekey_a);
  return k;
}

//input is 128 or 130 byte
function _utilPubkey2SecpFormat(utilPubKeyStr) {
  if(utilPubKeyStr.startsWith('0x')){
    utilPubKeyStr = utilPubKeyStr.slice(2);
  }
  utilPubKeyStr = '04' + utilPubKeyStr;
  return ethUtil.secp256k1.publicKeyConvert(new Buffer(utilPubKeyStr, 'hex'));
}

var tprivatekey = 'daa2fbee5ee569bc64842f5a386e7037612e0736b52e41749d52b616beaca65e';
var tAddress = '0xc29258c409380d34c9255406e8204212da552f92'

var TestPublicKey = getPubKeyStrByPrivateKey(tprivatekey);
var inputPubKey = _utilPubkey2SecpFormat(TestPublicKey);
var privateKey = new Buffer(tprivatekey, 'hex');

var retOTA = generateOTAPublicKey(inputPubKey,inputPubKey);
var otaPrivate = computeOTAPrivateKeys(privateKey, privateKey, retOTA.otaPubA1, retOTA.otaPubS1);

var checkPub = ethUtil.bufferToHex(retOTA.otaPubA1).slice(4)
var pubByOTAPrivate = getPubKeyStrByPrivateKey(ethUtil.bufferToHex(otaPrivate)).slice(2)
console.log(checkPub === pubByOTAPrivate ? "congratulation!" : "failed!!!");

function logSeperator(){
  console.log('   __    ___  @____  _     _     ___   _      ____  ____  _     ');
  console.log('  / /`  | |_)   / / \\ \\_/ | |   / / \\ \\ \\  / | |_  | |_  | |\\/| ');
  console.log('  \\_\\_, |_| \\  /_/_  |_|  |_|__ \\_\\_/  \\_\\/  |_|__ |_|__ |_|  | ');
  console.log('--------------------------------------------------------------------------')
}
