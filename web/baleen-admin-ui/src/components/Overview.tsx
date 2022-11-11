import { numberFormat } from 'utils/format';

type OverViewItemProps = {
    stats: number;
    title: string;
    icon?: string;
    iconClass?: string;
};
const OverView = ({ stats, title }: OverViewItemProps) => {
    return (
        <div className="d-flex p-2 border border rounded-1 overview">
            <div className="flex-grow-1">
                <span className="text-muted">{title}</span>
                <h4 className="mt-1 mb-1 fs-22">{numberFormat(stats)}</h4>
            </div>
        </div>
    );
};

export default OverView;
