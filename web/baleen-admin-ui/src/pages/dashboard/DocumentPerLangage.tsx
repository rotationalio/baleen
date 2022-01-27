import { Card } from 'react-bootstrap';
import Chart from 'react-apexcharts';
import { ApexOptions } from 'apexcharts';

const DocumentPerLangage = () => {
    const apexBarChartOpts: ApexOptions = {
        plotOptions: {
            pie: {
                donut: {
                    size: '70%',
                    labels: {
                        show: false,
                        name: {
                            show: true,
                            fontSize: '22px',
                            fontFamily: 'Helvetica, Arial, sans-serif',
                            fontWeight: 600,
                            color: undefined,
                            offsetY: -10,
                            formatter: function (val) {
                                return val;
                            },
                        },
                    },
                },
                expandOnClick: false,
            },
        },
        chart: {
            height: 291,
            type: 'donut',
        },
        legend: {
            show: true,
            position: 'bottom',
            horizontalAlign: 'center',
            itemMargin: {
                horizontal: 6,
                vertical: 3,
            },
        },
        labels: ['English', 'Spanish', 'French', 'Korean'],
        responsive: [
            {
                breakpoint: 480,
                options: {
                    legend: {
                        position: 'bottom',
                    },
                },
            },
        ],
        tooltip: {
            y: {
                formatter: (value: number) => {
                    return value + 'k';
                },
            },
        },
    };

    const apexBarChartData = [44, 55, 41, 17];

    return (
        <Card style={{ height: 450 }}>
            <Card.Body>
                <h5 className="card-title mt-0 mb-0 header-title">Documents per language</h5>

                <Chart
                    options={apexBarChartOpts}
                    series={apexBarChartData}
                    type="donut"
                    className="apex-charts mb-0 mt-4"
                    height={360}
                    dir="ltr"
                />
            </Card.Body>
        </Card>
    );
};

export default DocumentPerLangage;
