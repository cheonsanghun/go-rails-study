import type { SalesReport } from "@/types/domain";

// This file renders simple revenue data returned by the backend API.

type SalesReportTableProps = {
  reports: SalesReport[];
};

export function SalesReportTable({ reports }: SalesReportTableProps) {
  const totalUnits = reports.reduce((sum, report) => sum + report.units, 0);
  const totalRevenue = reports.reduce(
    (sum, report) => sum + report.netRevenueMinor,
    0,
  );
  const currency = reports[0]?.currency ?? "JPY";

  return (
    <section className="panel">
      <h2>수익 리포트</h2>
      <div className="meta-row">
        <span className="pill">Units: {totalUnits.toLocaleString()}</span>
        <span className="pill">
          Revenue: {formatRevenue(totalRevenue, currency)}
        </span>
      </div>

      {reports.length === 0 ? (
        <p className="muted">아직 수익 데이터가 없습니다.</p>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Period</th>
                <th>Track</th>
                <th>Store</th>
                <th>Country</th>
                <th>Units</th>
                <th>Revenue</th>
              </tr>
            </thead>
            <tbody>
              {reports.map((report) => (
                <tr key={report.id}>
                  <td>{report.period}</td>
                  <td>{report.trackTitle}</td>
                  <td>{report.storeName}</td>
                  <td>{report.country}</td>
                  <td>{report.units.toLocaleString()}</td>
                  <td>
                    {formatRevenue(report.netRevenueMinor, report.currency)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

function formatRevenue(minorAmount: number, currency: string) {
  const amount = currency === "JPY" ? minorAmount : minorAmount / 100;
  return new Intl.NumberFormat("ja-JP", {
    style: "currency",
    currency,
  }).format(amount);
}
