import { Col, Row } from 'react-bootstrap';
import DocumentPerLangage from './DocumentPerLangage';
import FeedsOverview from './FeedsOverview';
import Statistics from './Statistics';
import Timeline from './Timeline';

const Dashboard: React.FC = () => {
    return (
        <>
            <div className="page-title-box">
                <Row>
                    <Col>
                        <h4 className="page-title">Dashboard</h4>
                    </Col>
                </Row>
            </div>
            <Row>
                <Statistics />
            </Row>
            <Row>
                <Col md={6}>
                    <FeedsOverview />
                </Col>
                <Col md={6}>
                    <DocumentPerLangage />
                </Col>
            </Row>
            <Row>
                <Col>
                    <Timeline />
                </Col>
            </Row>
        </>
    );
};

export default Dashboard;
