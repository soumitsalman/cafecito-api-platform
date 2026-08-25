import { Badge } from "zudoku/ui/Badge";

/** Canonical public product status. Live APIs: Beans and Espresso only. */
export const CAFECITO_GATEWAY = "https://api.cafecito.tech";

export const PRODUCT_STATEMENTS = {
  beans:
    "The Beans API is a read-only publisher-content API for news, blogs, earnings and financial reports, litigation, official statements, research, technical documents, Sources, and Stories.",
  espresso:
    "The Espresso API is a read-only market and business intelligence API for searching Events, interpreting Signals, and tracing evidence and Sources.",
  cortado:
    "Cortado is a future social-media management API. It has no public REST or MCP surface.",
  latte:
    "Latte is a reserved future Cafecito product. It has no public REST or MCP surface.",
} as const;

export function ProductStatusTable() {
  return (
    <table>
      <thead>
        <tr>
          <th>Product</th>
          <th>Status</th>
          <th>Features</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>
            <strong>
              <a href="/products/beans">Beans</a>
            </strong>{" "}
            News and publisher content
          </td>
          <td>
            <Badge className="badge-live">Live</Badge>
          </td>
          <td>{PRODUCT_STATEMENTS.beans}</td>
        </tr>
        <tr>
          <td>
            <strong>
              <a href="/products/espresso">Espresso</a>
            </strong>{" "}
            Market and business intelligence
          </td>
          <td>
            <Badge className="badge-live">Live</Badge>
          </td>
          <td>{PRODUCT_STATEMENTS.espresso}</td>
        </tr>
        <tr>
          <td>
            <strong>
              <a href="/products/cortado">Cortado</a>
            </strong>{" "}
            Social media automation
          </td>
          <td>
            <Badge variant="secondary">Coming Soon</Badge>
          </td>
          <td>{PRODUCT_STATEMENTS.cortado}</td>
        </tr>
        <tr>
          <td>
            <strong>Latte</strong> Reserved product
          </td>
          <td>
            <Badge variant="secondary">Coming Soon</Badge>
          </td>
          <td>{PRODUCT_STATEMENTS.latte}</td>
        </tr>
      </tbody>
    </table>
  );
}
