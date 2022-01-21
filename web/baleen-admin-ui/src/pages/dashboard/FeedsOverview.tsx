import { Link } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import OverViewItem from 'components/OverViewItem';

const FeedsOverview = () => {
    return (
        <Card>
            <Card.Body className="p-0">
                <div className="p-3">
                    <span
                        className="float-end text-danger ms-2"
                        style={{ verticalAlign: 'middle', lineHeight: 'normal' }}>
                        3 Errored
                    </span>
                    <span className="float-end text-success" style={{ verticalAlign: 'middle', lineHeight: 'normal' }}>
                        25 Active
                    </span>
                    <h5 className="card-title header-title mb-0">Feeds Overview</h5>
                </div>

                <OverViewItem
                    link={'https://vietnamnews.vn'}
                    error={`403 Forbidden`}
                    icon={'rss'}
                    iconClass={'icon-xs icon-dual-danger'}
                />
                <OverViewItem
                    link={'https://themoscowtimes.com'}
                    error={`404 Not Found`}
                    icon={'rss'}
                    iconClass={'icon-xs icon-dual-danger'}
                />
                <OverViewItem
                    link={'dailynk.com/english'}
                    error={`Unexpected Error`}
                    icon={'rss'}
                    iconClass={'icon-xs icon-dual-danger'}
                />
                <OverViewItem link={'malay.com'} icon={'rss'} iconClass={'icon-xs icon-dual-success'} />
                <OverViewItem link={'dailynk.com/english'} icon={'rss'} iconClass={'icon-xs icon-dual-success'} />
                <OverViewItem link={'dailynk.com/english'} icon={'rss'} iconClass={'icon-xs icon-dual-success'} />
                <OverViewItem link={'dailynk.com/english'} icon={'rss'} iconClass={'icon-xs icon-dual-success'} />
                <OverViewItem link={'dailynk.com/english'} icon={'rss'} iconClass={'icon-xs icon-dual-success'} />

                <Link to="#" className="p-2 d-block text-end">
                    View All <i className="uil-arrow-right"></i>
                </Link>
            </Card.Body>
        </Card>
    );
};

export default FeedsOverview;
