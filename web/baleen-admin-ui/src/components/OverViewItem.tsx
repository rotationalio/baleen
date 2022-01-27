import FeatherIcon from 'feather-icons-react';
import classNames from 'classnames';
import { OverViewItemProps } from 'interfaces';

const OverViewItem = ({ link, error, icon, iconClass }: OverViewItemProps) => {
    return (
        <div className="justify-content-between align-items-center d-flex py-2 px-3 border-bottom">
            <p className="text-muted m-0 ms-1" style={{ lineHeight: 'normal', fontWeight: 300 }}>
                <FeatherIcon icon={icon} className={classNames('align-self-center', iconClass)} />
                <a href={link} target="_blank" rel="noreferrer" className="ms-1">
                    {link}
                </a>
            </p>
            <span className={'badge badge-soft-danger'}>{error}</span>
        </div>
    );
};

export default OverViewItem;
