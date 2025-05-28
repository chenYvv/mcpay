
// 导出 Excel文件
function exportExcel(data,name) {
  var sheet = XLSX.utils.aoa_to_sheet(data);
  openDownloadDialog(sheet2blob(sheet), name);
}

function IEVersion() {
  var userAgent = navigator.userAgent; //取得浏览器的userAgent字符串  
  var isIE = userAgent.indexOf("compatible") > -1 && userAgent.indexOf("MSIE") > -1; //判断是否IE<11浏览器  
  var isEdge = userAgent.indexOf("Edge") > -1 && !isIE; //判断是否IE的Edge浏览器  
  var isIE11 = userAgent.indexOf('Trident') > -1 && userAgent.indexOf("rv:11.0") > -1;
  if(isIE) {
      var reIE = new RegExp("MSIE (\\d+\\.\\d+);");
      reIE.test(userAgent);
      var fIEVersion = parseFloat(RegExp["$1"]);
      if(fIEVersion == 7) {
          return 7;
      } else if(fIEVersion == 8) {
          return 8;
      } else if(fIEVersion == 9) {
          return 9;
      } else if(fIEVersion == 10) {
          return 10;
      } else {
          return 6;//IE版本<=7
      }   
  } else if(isEdge) {
      return 'edge';//edge
  } else if(isIE11) {
      return 11; //IE11  
  }else{
      return -1;//不是ie浏览器
  }
}
/**
 * 通用的打开下载对话框方法，没有测试过具体兼容性
 * @param url 下载地址，也可以是一个blob对象，必选
 * @param saveName 保存文件名，可选
 */
function openDownloadDialog(url, saveName)
{
  if(typeof url == 'object' && url instanceof Blob)
  {
    url = URL.createObjectURL(url); // 创建blob地址
  }
  var aLink = document.createElement('a');
  aLink.href = url;
  aLink.download = saveName || ''; // HTML5新增的属性，指定保存文件名，可以不要后缀，注意，file:///模式下不会生效
  var event;
  
  if(window.MouseEvent){
    if (IEVersion()) {
      event = document.createEvent('MouseEvents');
      event.initMouseEvent('click', true, false, window, 0, 0, 0, 0, 0, false, false, false, false, 0, null);
    }else {
      event = new MouseEvent('click');
    }
  }
  else
  {
    event = document.createEvent('MouseEvents');
    event.initMouseEvent('click', true, false, window, 0, 0, 0, 0, 0, false, false, false, false, 0, null);
  }
  aLink.dispatchEvent(event);
}

// 将一个sheet转成最终的excel文件的blob对象，然后利用URL.createObjectURL下载
function sheet2blob(sheet, sheetName) {
  sheetName = sheetName || 'sheet1';
  var workbook = {
  SheetNames: [sheetName],
  Sheets: {}
  };
  workbook.Sheets[sheetName] = sheet;
  // 生成excel的配置项
  var wopts = {
  bookType: 'xlsx', // 要生成的文件类型
  bookSST: false, // 是否生成Shared String Table，官方解释是，如果开启生成速度会下降，但在低版本IOS设备上有更好的兼容性
  type: 'binary'
  };
  var wbout = XLSX.write(workbook, wopts);
  var blob = new Blob([s2ab(wbout)], {type:"application/octet-stream"});
  // 字符串转ArrayBuffer
  function s2ab(s) {
  var buf = new ArrayBuffer(s.length);
  var view = new Uint8Array(buf);
  for (var i=0; i!=s.length; ++i) view[i] = s.charCodeAt(i) & 0xFF;
  return buf;
  }
  return blob;
}