import { Card, Col, Row, Typography } from "antd";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { fetchDepartments } from "../api/catalog";
import { api } from "../api/client";
import { useI18n } from "../store/I18nContext";

export default function Home() {
  const [depts, setDepts] = useState<any[]>([]);
  const { t } = useI18n();
  useEffect(() => {
    fetchDepartments()
      .then(setDepts)
      .catch(() => setDepts([]));
  }, []);
  return (
    <div>
      <Typography.Title level={2}>{t("home.title")}</Typography.Title>
      <Typography.Paragraph>{t("home.subtitle")}</Typography.Paragraph>
      <Row gutter={[16, 16]}>
        {depts.map((d) => (
          <Col key={d.id} xs={24} sm={12} md={8}>
            <Link to={`/catalog?dept=${d.id}`} style={{ display: "block" }}>
              <Card
                hoverable
                title={d.name}
                cover={
                  d.image_url
                    ? (() => {
                        const origin = (() => {
                          try {
                            return new URL(api.defaults.baseURL as string)
                              .origin;
                          } catch {
                            return "";
                          }
                        })();
                        const img = !/^https?:/i.test(d.image_url)
                          ? origin + d.image_url
                          : d.image_url;
                        return <img src={img} alt={d.name} />;
                      })()
                    : undefined
                }
              >
                <p>{d.description}</p>
              </Card>
            </Link>
          </Col>
        ))}
      </Row>
    </div>
  );
}
