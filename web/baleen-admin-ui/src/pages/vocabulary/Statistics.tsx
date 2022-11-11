import OverView from 'components/Overview';
import { Col, Row } from 'react-bootstrap';
import { Vocabulary } from 'types';

type StatisticsProps = {
    data: Vocabulary;
};

function Statistics({ data }: StatisticsProps) {
    return (
        <>
            <Row className="mb-3">
                <Col>
                    <OverView stats={data.total_documents} title="Total documents" />
                </Col>
                <Col>
                    <OverView stats={data.total_words} title="Total words" />
                </Col>
                <Col>
                    <OverView stats={data.total_unique_words} title="Total unique words" />
                </Col>
                <Col>
                    <OverView stats={data.words_per_document} title="Words per document" />
                </Col>
            </Row>
        </>
    );
}

export default Statistics;
