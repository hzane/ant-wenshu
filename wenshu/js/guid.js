var createGuid = function () {
  return (((1 + Math.random()) * 65536) | 0).toString(16).substring(1);
};
function GetGUID(){
  var guid = createGuid() + createGuid() + '-' + createGuid() + '-' + createGuid() + createGuid() + '-' + createGuid() + createGuid() + createGuid();
  return guid;
};
