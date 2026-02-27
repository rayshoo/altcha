(function () {
  "use strict";

  let chart = null;

  const rangeSelect = document.getElementById("range-select");
  const customRange = document.getElementById("custom-range");
  const dateFrom = document.getElementById("date-from");
  const dateTo = document.getElementById("date-to");
  const btnApply = document.getElementById("btn-apply");

  rangeSelect.addEventListener("change", function () {
    if (this.value === "custom") {
      customRange.classList.remove("hidden");
    } else {
      customRange.classList.add("hidden");
      refreshDashboard();
    }
  });

  btnApply.addEventListener("click", function () {
    refreshDashboard();
  });

  function getTimeRange() {
    const days = rangeSelect.value;
    if (days === "custom") {
      return { from: dateFrom.value, to: dateTo.value };
    }
    const to = new Date();
    const from = new Date();
    from.setDate(from.getDate() - parseInt(days));
    return {
      from: formatDate(from),
      to: formatDate(to),
    };
  }

  function formatDate(d) {
    return d.toISOString().slice(0, 10);
  }

  async function fetchJSON(url) {
    const resp = await fetch(url);
    if (resp.status === 401) {
      window.location.reload();
      throw new Error("Unauthorized");
    }
    if (!resp.ok) throw new Error("HTTP " + resp.status);
    return resp.json();
  }

  async function refreshDashboard() {
    const range = getTimeRange();
    const qs = "from=" + range.from + "&to=" + range.to;

    try {
      const [summary, timeseries, locations] = await Promise.all([
        fetchJSON("/api/summary?" + qs),
        fetchJSON("/api/timeseries?" + qs),
        fetchJSON("/api/locations?" + qs),
      ]);

      renderKPICards(summary);
      renderChart(timeseries);
      renderLocations(locations);
    } catch (err) {
      console.error("Dashboard refresh error:", err);
    }
  }

  function renderKPICards(s) {
    const grid = document.getElementById("kpi-grid");
    const cards = [
      { label: "Challenges", value: fmtNum(s.challenges), cls: "" },
      { label: "Verified", value: fmtNum(s.verified), cls: "success" },
      { label: "Failed", value: fmtNum(s.failed), cls: "error" },
      {
        label: "Avg Latency",
        value: s.avg_latency_ms.toFixed(1) + " ms",
        cls: "",
      },
      { label: "4XX Errors", value: fmtNum(s.errors_4xx), cls: "error" },
      { label: "5XX Errors", value: fmtNum(s.errors_5xx), cls: "error" },
      { label: "Requests", value: fmtNum(s.total_requests), cls: "" },
    ];

    grid.innerHTML = cards
      .map(
        (c) =>
          '<div class="kpi-card ' +
          c.cls +
          '">' +
          '<div class="kpi-value">' +
          c.value +
          "</div>" +
          '<div class="kpi-label">' +
          c.label +
          "</div>" +
          "</div>"
      )
      .join("");
  }

  function renderChart(data) {
    const ctx = document.getElementById("main-chart").getContext("2d");

    if (chart) {
      chart.destroy();
    }

    const labels = data.map(function (d) {
      return d.date;
    });

    chart = new Chart(ctx, {
      type: "bar",
      data: {
        labels: labels,
        datasets: [
          {
            type: "bar",
            label: "Challenges",
            data: data.map(function (d) {
              return d.challenges;
            }),
            backgroundColor: "rgba(59, 130, 246, 0.7)",
            yAxisID: "y",
            order: 2,
          },
          {
            type: "bar",
            label: "Verified",
            data: data.map(function (d) {
              return d.verified;
            }),
            backgroundColor: "rgba(43, 196, 172, 0.7)",
            yAxisID: "y",
            order: 2,
          },
          {
            type: "bar",
            label: "Failed",
            data: data.map(function (d) {
              return d.failed;
            }),
            backgroundColor: "rgba(231, 76, 60, 0.7)",
            yAxisID: "y",
            order: 2,
          },
          {
            type: "line",
            label: "Avg Latency (ms)",
            data: data.map(function (d) {
              return d.avg_latency;
            }),
            borderColor: "#f39c12",
            backgroundColor: "rgba(243, 156, 18, 0.1)",
            yAxisID: "y1",
            tension: 0.3,
            pointRadius: 3,
            order: 1,
          },
        ],
      },
      options: {
        responsive: true,
        interaction: { mode: "index", intersect: false },
        scales: {
          y: {
            type: "linear",
            position: "left",
            beginAtZero: true,
            title: { display: true, text: "Count" },
          },
          y1: {
            type: "linear",
            position: "right",
            beginAtZero: true,
            title: { display: true, text: "Latency (ms)" },
            grid: { drawOnChartArea: false },
          },
        },
      },
    });
  }

  // Country code → flag emoji
  var FLAG_OFFSET = 0x1f1e6 - 65;
  function countryFlag(code) {
    if (!code || code.length !== 2 || code === "Unknown") return "";
    var a = code.charCodeAt(0);
    var b = code.charCodeAt(1);
    return (
      String.fromCodePoint(a + FLAG_OFFSET) +
      String.fromCodePoint(b + FLAG_OFFSET)
    );
  }

  var CONTINENT_NAMES = {
    AF: "Africa",
    AN: "Antarctica",
    AS: "Asia",
    EU: "Europe",
    NA: "North America",
    OC: "Oceania",
    SA: "South America",
    Unknown: "Unknown",
  };

  function renderLocations(data) {
    var container = document.getElementById("locations-container");
    var section = document.getElementById("locations-section");

    if (!data || data.length === 0) {
      section.classList.add("hidden");
      return;
    }
    section.classList.remove("hidden");

    // Check if all locations are unknown (GeoIP not configured)
    var allUnknown = data.every(function (d) {
      return d.continent === "Unknown" && d.country === "Unknown";
    });
    if (allUnknown) {
      section.classList.add("hidden");
      return;
    }

    // Group by continent
    var groups = {};
    data.forEach(function (entry) {
      var continent = entry.continent || "Unknown";
      if (!groups[continent]) groups[continent] = [];
      groups[continent].push(entry);
    });

    var html = "";
    var continentOrder = ["NA", "EU", "AS", "SA", "AF", "OC", "AN", "Unknown"];
    continentOrder.forEach(function (cont) {
      var entries = groups[cont];
      if (!entries) return;

      var name = CONTINENT_NAMES[cont] || cont;
      html += '<div class="continent-group"><h3>' + name + "</h3>";
      entries.forEach(function (e) {
        html +=
          '<div class="location-row">' +
          '<span class="location-flag">' +
          countryFlag(e.country) +
          "</span>" +
          '<span class="location-country">' +
          e.country +
          "</span>" +
          '<div class="location-bar-bg"><div class="location-bar" style="width:' +
          e.percent.toFixed(1) +
          '%"></div></div>' +
          '<span class="location-count">' +
          fmtNum(e.count) +
          "</span>" +
          '<span class="location-pct">' +
          e.percent.toFixed(1) +
          "%</span>" +
          "</div>";
      });
      html += "</div>";
    });

    container.innerHTML = html;
  }

  function fmtNum(n) {
    if (n == null) return "0";
    return n.toLocaleString();
  }

  // Initial load
  refreshDashboard();
})();
