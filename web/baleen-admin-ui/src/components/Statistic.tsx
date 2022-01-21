import React from 'react';
import { Card } from 'react-bootstrap';
import FeatherIcon from 'feather-icons-react';
import { numberFormat } from 'utils/format';

interface StatisticsChartWidgetProps {
    title: string;
    stats: number | string | bigint;
    iconProps?: {
        icon?: string;
        className?: string;
    };
}

const StatisticsChartWidget = ({ title, stats, iconProps }: StatisticsChartWidgetProps) => {
    const _stats = React.useMemo(
        () => (typeof stats === 'number' || typeof stats === 'bigint' ? numberFormat(stats) : stats),
        [stats]
    );

    return (
        <Card className="">
            <Card.Body>
                <div className="d-flex">
                    <div className="flex-grow-1">
                        <span className="text-muted text-uppercase fs-8 fw-bold">{title}</span>
                        <h3 className="mb-0">{_stats}</h3>
                    </div>
                    <div className="align-self-center flex-shrink-0">
                        <FeatherIcon size={25} {...iconProps} />
                    </div>
                </div>
            </Card.Body>
        </Card>
    );
};

export default StatisticsChartWidget;
