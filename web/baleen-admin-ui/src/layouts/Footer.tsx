import { Row, Col } from 'react-bootstrap';

const Footer = () => {
    const currentYear = new Date().getFullYear();

    return (
        <footer className="footer">
            <div className="container-fluid">
                <Row>
                    <Col sm={6}>{currentYear} &copy; Rotational Labs, LLC, All Rights Reserved</Col>
                </Row>
            </div>
        </footer>
    );
};

export default Footer;
