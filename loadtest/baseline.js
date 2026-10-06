// Baseline load test. Run it BEFORE you change anything and keep the numbers —
// every improvement you claim in a pull request is measured against this run.
//
//   k6 run loadtest/baseline.js
//
// Install k6: https://grafana.com/docs/k6/latest/set-up/install-k6/
import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 10,
  duration: "30s",
  thresholds: {
    // The scaffold FAILS this threshold on purpose. Making it pass is the work.
    http_req_duration: ["p(95)<100"],
  },
};

const BASE = __ENV.BASE_URL || "http://localhost:8080";

export default function () {
  const id = Math.floor(Math.random() * 1000) + 1;
  const res = http.get(`${BASE}/products/${id}`);
  check(res, { "status 200": (r) => r.status === 200 });
}
