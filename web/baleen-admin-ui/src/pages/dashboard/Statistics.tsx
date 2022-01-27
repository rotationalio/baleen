import Statistic from 'components/Statistic';
import { Row, Col } from 'react-bootstrap';

// components

const Statistics = () => {
    return (
        <>
            <Row>
                <Col>
                    <Statistic
                        iconProps={{ icon: 'rss', className: 'icon-dual-success' }}
                        title="Total Feeds"
                        stats={3000}
                    />
                </Col>
                <Col>
                    <Statistic
                        iconProps={{ icon: 'arrow-down-circle', className: 'icon-dual-success' }}
                        title="Total File Size"
                        stats={`30 Go`}
                    />
                </Col>

                <Col>
                    <Statistic
                        iconProps={{ icon: 'file-text', className: 'icon-dual-info' }}
                        title="Total Documents Processed"
                        stats={3000}
                    />
                </Col>

                <Col>
                    <Statistic
                        iconProps={{ icon: 'search', className: 'icon-dual-warning' }}
                        title="Total Words"
                        stats={30000}
                    />
                </Col>
                <Col>
                    <Statistic
                        iconProps={{ icon: 'check', className: 'icon-dual-danger' }}
                        title="Total Unique Words"
                        stats={300}
                    />
                </Col>
            </Row>
        </>
    );
};

export default Statistics;
