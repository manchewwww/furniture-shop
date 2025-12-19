import { Card, Col, Row, Typography } from "antd";
import { useEffect, useState } from "react";
import { fetchDepartments } from "../api/catalog";
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
            <Card title={d.name}>
              <p>{d.description}</p>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}
