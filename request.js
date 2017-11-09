"use strict";

var system = require('system');
var server = '';

if (system.args.length === 1) {
    console.log('need url, like: http://www.taobao.com.');
} else {
    system.args.forEach(function (arg, i) {
        server = arg;
    });
}

var webPage = require('webpage');
var page = webPage.create();

// 如果页面采用了懒加载：1，增大尺寸，使底部可见
page.viewportSize = {
    width: 4800,
    height: 8000
};
page.settings.userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36";
page.settings.encoding = "utf8"; //"GBK";

page.open(server, function (status) {

    // 异步拉取数据的网站，需要等待一断时间后再去操作dom，时间长短由网站大小决定
    setTimeout(function () {
        // var result = page.evaluate(function () {
        // 如果页面采用了懒加载：2，滚动到底部
        // window.scrollTo(0, 10000);

        // 点击‘加载更多’
        // document.getElementsByClassName('get-more-line').click();
        // 
        //     return document.title;
        // });
        // console.log(result);

        console.log(page.content);

        //生成当前页面截图
        // page.render("snapshot.png");

        phantom.exit();
    }, 15000);
});