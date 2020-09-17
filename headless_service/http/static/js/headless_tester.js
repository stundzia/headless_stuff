$(document).ready(function () {
    var getResultsButton = $("button#get_results");
    var targetUrl = $("input#target_url");
    var waitTime = $("input#wait_time");
    var engine = $("select#engine");
    var uaType = $("select#user_agent_type");
    var ua = $("input#user_agent");
    var proxy = $("input#proxy_url");
    var headless = $("input#headless_mode");

    var onchangeEngine = function (self) {
        if (self.val() === 'surf') {
            headless.prop("checked", true);
            headless.prop("disabled", true);
        } else {
            headless.prop("disabled", false);
        }
    };

    onchangeEngine(engine);
    engine.change(function () {
        onchangeEngine($(this))
    });

    getResultsButton.click(function () {
        var request_time = new Date();
        var payload = {
            "url": targetUrl.val(),
            "wait_time": parseFloat(waitTime.val()),
            "engine": engine.val(),
            "user_agent_type": uaType.val(),
            "user_agent": ua.val(),
            "proxy": proxy.val(),
            "headless": headless.prop("checked"),
        };
        console.log("payload: ", payload);
        targetUrl.prop("disabled", true);
        getResultsButton.prop("disabled", true);
        $.ajax({
            type: "POST",
            url: "render_get",
            data: JSON.stringify(payload),
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            success: function (data) {
                console.log(data);
                $("textarea#response").val(JSON.stringify(data, undefined, 2));
                targetUrl.prop("disabled", false);
                getResultsButton.prop("disabled", false);
                var showHtml = $("input#open_html").prop("checked");
                if (showHtml === true) {
                    var rr = data['content'];
                    var iframe_div = $("div.embed-responsive.embed-responsive-16by9");
                    $("iframe#iframe_html_res").remove();
                    // TODO: will be needed if img rendering will be implemented.
                    // if (splashType.prop("checked") === true) {
                    //     var scrn = '<img src="data:image/png;base64,' + rr + '" />'
                    // }
                    var iframe_el = $("<iframe id=\"iframe_html_res\" class=\"embed-responsive-item\"></iframe>");
                    iframe_div.append(iframe_el);
                    var iframe_doc = document.getElementById('iframe_html_res').contentWindow.document;
                    iframe_doc.open();
                    iframe_doc.clear();
                    iframe_doc.write(rr);
                    iframe_doc.close();
                    $("div#html_res_row").show();
                }
                var response_time = new Date();
                var response_dur = (response_time - request_time) / 1000;
                $("label#response_label")[0].innerText = "Response received in " + response_dur + " seconds:";
            },
            failure: function (errMsg) {
                console.log(errMsg);
                $("textarea#response").val(errMsg);
                targetUrl.prop("disabled", false);
                getResultsButton.prop("disabled", false);
            }
        });
    })
});
