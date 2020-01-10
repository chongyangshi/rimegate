REACHABILITY_CHECK_INVERVAL = 30;

var apiBase = "https://rimegate.dogelink.com";
var grafanaUsername = "";
var grafanaPassword = "";

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
    request.open("GET", apiBase + "/ping", false);
    request.send();

    if (request.status !== 200) {
        console.log("Connectivity ping failed: " + request.responseText);
        document.getElementById("STATUS").innerText = "Disconnected"
        document.getElementsByClassName("header-right")[0].classList.remove("status-connected")
        document.getElementsByClassName("header-right")[0].classList.add("status-disconnected")
        return
    }

    document.getElementById("STATUS").innerText = "Connected"
    document.getElementsByClassName("header-right")[0].classList.remove("status-disconnected")
    document.getElementsByClassName("header-right")[0].classList.add("status-connected")
}

function listDashboards() {
    var request = new XMLHttpRequest();
    request.open("POST", apiBase + "/dashboards", false);
    request.setRequestHeader("Content-Type", "application/json");
    request.send(JSON.stringify({ "grafana_username": grafanaUsername, "grafana_password": grafanaPassword }));

    var response = JSON.parse(request.responseText);

    if (request.status !== 200) {
        console.log("Dashboard list failed (" + response.code + "): " + response.message);
        document.getElementsByClassName("error-text")[0].firstElementChild.innerText = "Error " + response.code + ": " + response.message;
        document.getElementsByClassName("error-text")[0].style.display = "block";
        return false
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
                dashboardButton.appendChild(button);
            }
        }
    }

    document.getElementById("INSTRUCTION").innerText = "Select Dashboard"
    document.getElementsByClassName("dashboard-selection")[0].style.display = "block";

    return true
}

function login() {
    grafanaUsername = document.getElementsByName("grafana_username")[0].value;
    grafanaPassword = document.getElementsByName("grafana_password")[0].value;

    var result = listDashboards();
    if (result === true) {
        document.getElementsByClassName("login")[0].style.display = "none";
    }
}

window.onload = init;
