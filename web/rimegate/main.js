REACHABILITY_CHECK_INVERVAL = 30;
DASHBOARD_REFRESH_INTERVAL = 30;

var apiBase = "API_BASE";
var grafanaUsername = "";
var grafanaPassword = "";
var connected = false;
var dashboardURL = "";
var dashboardTitle = "";
var autoFitPanel = true;

if (window.location.hostname == undefined || window.location.hostname == "") {
    // For local running, served by browser from file.
    apiBase = "http://127.0.0.1:18080";
}

function init() {
    ping();
    setInterval(ping, REACHABILITY_CHECK_INVERVAL * 1000);

    document.getElementById("LOGIN").addEventListener("click", login);
    document.getElementsByName("grafana_password")[0].addEventListener("keyup", function () {
        if (event.keyCode === 13) {
            event.preventDefault();
            login();
        }
    });
}

function ping() {
    var request = new XMLHttpRequest();
    request.onreadystatechange = done;
    request.open("GET", apiBase + "/ping");
    request.send();

    function done() {
        if (request.readyState !== XMLHttpRequest.DONE) {
            return
        }

        if (request.status !== 200) {
            console.log("Connectivity ping failed: " + request.responseText);
            document.getElementById("STATUS").innerText = "Disconnected";
            document.getElementsByClassName("header-right")[0].classList.remove("status-connected");
            document.getElementsByClassName("header-right")[0].classList.add("status-disconnected");
            connected = false
            return
        }

        document.getElementById("STATUS").innerText = "Connected";
        document.getElementsByClassName("header-right")[0].classList.remove("status-disconnected");
        document.getElementsByClassName("header-right")[0].classList.add("status-connected");
        connected = true;
    }

}

function listDashboards() {
    var request = new XMLHttpRequest();
    request.onreadystatechange = done;
    request.open("POST", apiBase + "/dashboards");
    request.setRequestHeader("Content-Type", "application/json");
    request.send(JSON.stringify({ "grafana_username": grafanaUsername, "grafana_password": grafanaPassword }));

    function done() {
        if (request.readyState !== XMLHttpRequest.DONE) {
            return
        }

        var response = JSON.parse(request.responseText);

        if (request.status !== 200) {
            console.log("Dashboard list failed (" + response.code + "): " + response.message);
            document.getElementsByClassName("error-text")[0].firstElementChild.innerText = "Error " + response.code + ": " + response.message;
            document.getElementsByClassName("error-text")[0].style.display = "block";
            return
        }

        document.getElementsByClassName("error-text")[0].firstElementChild.innerText = "";
        document.getElementsByClassName("error-text")[0].style.display = "none";

        // Clear the table if required
        var table = document.getElementById("DASHBOARDS");
        var rows = table.getElementsByTagName('tbody')[0];
        if (rows.length > 0) {
            for (var child in rows.children) {
                rows.removeChild(child);
            }
        }

        // Populate the table
        for (var folderIndex in response.dashboards) {
            if (response.dashboards.hasOwnProperty(folderIndex)) {
                var folder = response.dashboards[folderIndex];
                for (var i = 0; i < folder.length; i++) {
                    var row = rows.insertRow();
                    var folderName = row.insertCell(0);
                    var dashboardName = row.insertCell(1);
                    var dashboardButton = row.insertCell(2);

                    var folderNameText = folder[i].folderTitle;
                    if (folderNameText === "") {
                        folderNameText = "Default";
                    }

                    folderName.appendChild(document.createTextNode(folderNameText));
                    dashboardName.appendChild(document.createTextNode(folder[i].title));

                    var button = document.createElement("BUTTON");
                    button.innerText = "Launch";
                    button.setAttribute("data-dashboard-url", folder[i].url);
                    button.setAttribute("data-dashboard-title", folder[i].title);
                    button.addEventListener("click", launchDashboard);
                    dashboardButton.appendChild(button);
                }
            }
        }

        document.getElementById("INSTRUCTION").innerText = "Select Dashboard"
        document.getElementsByClassName("dashboard-selection")[0].style.display = "block";
        document.getElementsByClassName("login")[0].style.display = "none";
    }
}

function launchDashboard(e) {
    dashboardURL = e.target.getAttribute("data-dashboard-url");
    dashboardTitle = e.target.getAttribute("data-dashboard-title");
    document.getElementsByClassName("dashboard-selection")[0].style.display = "none";
    refreshDashboard();
    setInterval(refreshDashboard, DASHBOARD_REFRESH_INTERVAL * 1000);
}

function enableAutoFitPanelCheckbox() {
    document.getElementsByClassName("fit-panel")[0].style.display = "inline";
    document.getElementById("PANELBOX").addEventListener('change', function (e) {
        autoFitPanel = e.target.checked;
        refreshDashboard();
    });
}

function login() {
    grafanaUsername = document.getElementsByName("grafana_username")[0].value;
    grafanaPassword = document.getElementsByName("grafana_password")[0].value;

    listDashboards();
}

function refreshDashboard() {
    var width = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
    var height = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);

    document.getElementById("TICKER").innerText = "Loading..."

    var request = new XMLHttpRequest();
    request.onreadystatechange = done;
    request.open("POST", apiBase + "/render");
    request.setRequestHeader("Content-Type", "application/json");
    request.send(JSON.stringify({
        "dashboard_url": dashboardURL,
        "width": width,
        "height": height,
        "grafana_username": grafanaUsername,
        "grafana_password": grafanaPassword,
        "auto_fit_panel": autoFitPanel,
    }));

    function done() {
        if (request.readyState !== XMLHttpRequest.DONE) {
            return
        }

        var response = JSON.parse(request.responseText);

        time = parseISODate(response.rendered_time);

        document.getElementById("RENDER").src = "data:image/png;base64," + response.payload;
        document.getElementById("RENDER").style.display = "block";
        document.getElementById("INSTRUCTION").innerText = dashboardTitle;
        document.getElementById("TICKER").innerText = String(time.getHours()).padStart(2, "0") + ":" + String(time.getMinutes()).padStart(2, "0") + ":" + String(time.getSeconds()).padStart(2, "0");
        console.log("Rendered payload " + response.payload.length + " for dashboard " + dashboardTitle + " at " + response.rendered_time);

        enableAutoFitPanelCheckbox();
    }
}

// Borrowed from https://stackoverflow.com/a/11318669 comment by https://stackoverflow.com/users/528263/crashthatch
var dateRegex = /^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2}):(\d{2})\.?(\d{3})?(?:(?:([+-]\d{2}):?(\d{2}))|Z)?$/;
function parseISODate(d) {
    var m = dateRegex.exec(d);
    //milliseconds are optional.
    if (m[7] === undefined) {
        m[7] = 0;
    }

    //if timezone is undefined, it must be Z or nothing (otherwise the group would have captured).
    if (m[8] === undefined && m[9] === undefined) {
        //Use UTC.
        m[8] = 0;
        m[9] = 0;
    }

    var year = +m[1];
    var month = +m[2];
    var day = +m[3];
    var hour = +m[4];
    var minute = +m[5];
    var second = +m[6];
    var msec = +m[7];
    var tzHour = +m[8];
    var tzMin = +m[9];
    var tzOffset = tzHour * 60 + tzMin;

    //console.log(year+', '+(month - 1)+', '+day+', '+hour+', '+(minute - tzOffset)+', '+second+', '+msec);

    return new Date(year, month - 1, day, hour, minute - tzOffset, second, msec);
}

window.onload = init;
