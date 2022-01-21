import { Card } from 'react-bootstrap';
import Chart from 'react-apexcharts';
import { ApexOptions } from 'apexcharts';

const Timeline = () => {
    const apexBarChartOpts: ApexOptions = {
        chart: {
            height: 349,
            type: 'bar',
            stacked: true,
            toolbar: {
                show: false,
            },
        },
        plotOptions: {
            bar: {
                horizontal: false,
                columnWidth: '45%',
            },
        },
        dataLabels: {
            enabled: false,
        },
        stroke: {
            show: true,
            width: 2,
            colors: ['transparent'],
        },
        xaxis: {
            categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
            axisBorder: {
                show: false,
            },
        },
        legend: {
            show: false,
        },
        grid: {
            row: {
                colors: ['transparent', 'transparent'], // takes an array which will be repeated on columns
                opacity: 0.2,
            },
            borderColor: '#f3f4f7',
        },
        tooltip: {
            theme: 'dark',
            x: {
                show: false,
            },
            y: {
                formatter: function (val) {
                    return val + ' feeds';
                },
            },
        },
    };

    const apexBarChartData = [
        {
            name: 'Document processed',
            data: [35, 44, 55, 57, 56, 61],
        },
        {
            name: 'Feeds',
            data: [52, 76, 85, 101, 98, 87],
        },
    ];

    return (
        <Card>
            <Card.Body className="pb-0">
                <h5 className="card-title header-title">Timeline (Documents Count)</h5>

                <Chart
                    options={apexBarChartOpts}
                    series={apexBarChartData}
                    type="bar"
                    className="apex-charts mt-3"
                    height={349}
                    dir="ltr"
                />
            </Card.Body>
        </Card>
    );
};

export default Timeline;
