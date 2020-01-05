var API_BASE = "https://rimegate.dogelink.com"

if (window.location.hostname == undefined || window.location.hostname == "") {
    // For local running, served by browser from file.
    API_BASE = "http://127.0.0.1:18080";
}

function listDashboards() {
    
}