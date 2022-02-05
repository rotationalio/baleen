import { Row, Col, Card } from 'react-bootstrap';

import Table from 'components/Table';
import { DatatableProps, Vocabulary } from 'types';
import Statistics from './Statistics';

const columns = [
    {
        Header: 'Word',
        accessor: 'word',
        sort: true,
    },
    {
        Header: 'Count',
        accessor: 'count',
        sort: true,
    },
    {
        Header: 'Percentage',
        accessor: 'percentage',
        sort: false,
    },
];

const sizePerPageList = [
    {
        text: '5',
        value: 5,
    },
    {
        text: '10',
        value: 10,
    },
    {
        text: '25',
        value: 25,
    },
];

const Datatable: React.FC<DatatableProps> = ({ title, data }) => {
    const formatData = (data: Vocabulary): Record<string, any>[] => {
        return Object.entries(data.most_common_words).map(([k, v]) => ({ word: k, ...v }));
    };

    return (
        <>
            <Row>
                <Col>
                    <Card>
                        <Card.Body>
                            <h4 className="header-title mb-3">{title}</h4>
                            <Statistics data={data} />
                            <Table
                                columns={columns}
                                data={formatData(data) || []}
                                pageSize={5}
                                sizePerPageList={sizePerPageList}
                                isSortable={true}
                                pagination={true}
                            />
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </>
    );
};

export default Datatable;
