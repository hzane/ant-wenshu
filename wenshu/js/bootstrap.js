// 其中某些脚本依赖 window, document and setTimeout
var document = {};

var window = {
    innerWidth: 800,
    innerHeight: 600,
    screenX: 0,
    screenY: 0,
    location: "",
    atob: function (str) {
        return Base64_Zip.atob(str);
    },
    alert: function (str) {
        return str;
    },
    screen: {
        width: 1600,
        height: 1200,
    }
};
String.prototype.fontcolor = function (arg) {
    return arg;
};

function atob(str) {
    return Base64_Zip.atob(str);
}

function alert(str) {
    return str;
}

// " setTimeout('com.str._KEY=\"QQCbNrOvGfcZ6BdoGMVi0o/nhUSsc5Re\";',8000*Math.random());setTimeout('var a=0;while(1){a++;}',10)"
// setTimeout里面出现了一个循环, after 20190415
function setTimeout(a, b) {
    if (a.indexOf("com.str._KEY") >= 0) {
        eval(a)
    }

    // com.str._KEY = "QQCbNrOvGfcZ6BdoGMVi0o/nhUSsc5Re";
}

function CrashRunEval(runeval) {
    eval(unzip(runeval));
    return com.str._KEY
}

function CrashDOCID(did) {
    return com.str.Decrypt(unzip(did));
}

var createGuid = function () {
    return (((1 + Math.random()) * 65536) | 0).toString(16).substring(1);
}

function CrashGUID() {
    var guid = createGuid() + createGuid() + '-' + createGuid() + '-' + createGuid() + createGuid() + '-' + createGuid() + createGuid() + createGuid();
    return guid;
}
