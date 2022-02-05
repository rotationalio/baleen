import React, { useContext } from 'react';
import { Card, Col, Row, useAccordionButton, AccordionContext, Accordion } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import classNames from 'classnames';

interface ContentType {
    id: number;
    title: string;
    icon?: string;
    text: string;
}

interface CustomToggleProps {
    children: any;
    eventKey: string;
    containerClass: string;
    linkClass: string;
    callback?: any;
}

const CustomToggle = ({ children, eventKey, containerClass, linkClass, callback }: CustomToggleProps) => {
    const { activeEventKey } = useContext(AccordionContext);

    const decoratedOnClick = useAccordionButton(eventKey, () => callback && callback(eventKey));

    const isCurrentEventKey = activeEventKey === eventKey;

    return (
        <Link
            to="#"
            className={classNames(linkClass, {
                collapsed: !isCurrentEventKey,
            })}
            onClick={decoratedOnClick}>
            <Card.Header>
                <h5 className={containerClass}>{children}</h5>
            </Card.Header>
        </Link>
    );
};

const CustomAccordion1 = ({ item, index, length }: { item: ContentType; index: number; length: number }) => {
    return (
        <Card className={classNames('shadow-none', 'border', index + 1 === length ? 'mb-0' : 'mb-1')}>
            <CustomToggle eventKey={String(index)} containerClass="m-0 fs-16" linkClass="text-dark">
                {item.title}
                <i className="uil uil-angle-down float-end accordion-arrow"></i>
            </CustomToggle>
            <Accordion.Collapse eventKey={String(index)}>
                <Card.Body className="text-muted">{item.text}</Card.Body>
            </Accordion.Collapse>
        </Card>
    );
};

const CustomAccordion2 = ({ item, index, length }: { item: ContentType; index: number; length: number }) => {
    return (
        <Card className={classNames(index + 1 === length ? 'mb-0' : 'mb-1')}>
            <CustomToggle eventKey={String(index)} containerClass="m-0 fs-14" linkClass="text-dark">
                {item.title}
                <i className="uil uil-question-circle h3 text-primary icon"></i>
            </CustomToggle>
            <Accordion.Collapse eventKey={String(index)}>
                <Card.Body className="text-muted">{item.text}</Card.Body>
            </Accordion.Collapse>
        </Card>
    );
};

const Accordions = () => {
    const accordianContent: ContentType[] = [
        {
            id: 2,
            title: ' What is Lorem Ipsum?',
            text: `Anim pariatur cliche reprehenderit, enim eiusmod high life accusamus terry
                    richardson ad squid. 3 wolf moon officia aute, non cupidatat skateboard dolor brunch.Food truck quinoa
                    nesciunt laborum eiusmod.Brunch 3 wolf moon tempor, sunt aliqua
                    put a bird on it squid single- origin coffee nulla assumenda
                    shoreditch et.Nihil anim keffiyeh helvetica, craft beer labore
                    wes anderson cred nesciunt sapiente ea proident.Ad vegan
                    excepteur butcher vice lomo.Leggings occaecat craft beer
                    farm - to - table, raw denim aesthetic synth nesciunt you probably
                    haven't heard of them accusamus labore sustainable VHS.`,
        },
        {
            id: 3,
            title: 'Why do we use it?',
            text: `Anim pariatur cliche reprehenderit, enim eiusmod high life accusamus terry
                    richardson ad squid. 3 wolf moon officia aute, non cupidatat skateboard dolor brunch.Food truck quinoa
                    nesciunt laborum eiusmod.Brunch 3 wolf moon tempor, sunt aliqua
                    put a bird on it squid single- origin coffee nulla assumenda
                    shoreditch et.Nihil anim keffiyeh helvetica, craft beer labore
                    wes anderson cred nesciunt sapiente ea proident.Ad vegan
                    excepteur butcher vice lomo.Leggings occaecat craft beer
                    farm - to - table, raw denim aesthetic synth nesciunt you probably
                    haven't heard of them accusamus labore sustainable VHS.`,
        },
        {
            id: 4,
            title: 'Where does it come from?',
            text: `Anim pariatur cliche reprehenderit, enim eiusmod high life accusamus terry
                    richardson ad squid. 3 wolf moon officia aute, non cupidatat skateboard dolor brunch.Food truck quinoa
                    nesciunt laborum eiusmod.Brunch 3 wolf moon tempor, sunt aliqua
                    put a bird on it squid single- origin coffee nulla assumenda
                    shoreditch et.Nihil anim keffiyeh helvetica, craft beer labore
                    wes anderson cred nesciunt sapiente ea proident.Ad vegan
                    excepteur butcher vice lomo.Leggings occaecat craft beer
                    farm - to - table, raw denim aesthetic synth nesciunt you probably
                    haven't heard of them accusamus labore sustainable VHS.`,
        },
    ];

    const accordianContent2: ContentType[] = [
        {
            id: 1,
            title: 'What is Anim text here?',
            text: ` Anim pariatur cliche reprehenderit, enim eiusmod high life
                    accusamus terry richardson ad squid. 3 wolf moon officia aute,
                    non cupidatat skateboard dolor brunch. Food truck quinoa
                    nesciunt laborum eiusmod. Brunch 3 wolf moon tempor, sunt aliqua
                    put a bird on it squid single-origin coffee nulla assumenda
                    shoreditch et. Nihil anim keffiyeh helvetica craft beer. Ad
                    vegan excepteur butcher vice lomo. Leggings occaecat craft beer
                    farm-to-table, raw denim aesthetic synth nesciunt you probably
                    haven't heard of them accusamus labore sustainable VHS.`,
        },
        {
            id: 2,
            title: 'Why use Anim text here?',
            text: `Anim pariatur cliche reprehenderit, enim eiusmod high life accusamus terry
                    richardson ad squid. 3 wolf moon officia aute, non cupidatat skateboard dolor brunch.Food truck quinoa
                    nesciunt laborum eiusmod.Brunch 3 wolf moon tempor, sunt aliqua
                    put a bird on it squid single- origin coffee nulla assumenda
                    shoreditch et.Nihil anim keffiyeh helvetica, craft beer labore
                    wes anderson cred nesciunt sapiente ea proident.Ad vegan
                    excepteur butcher vice lomo.Leggings occaecat craft beer
                    farm - to - table, raw denim aesthetic synth nesciunt you probably
                    haven't heard of them accusamus labore sustainable VHS.`,
        },
        {
            id: 3,
            title: 'How many variations exist?',
            text: `Anim pariatur cliche reprehenderit, enim eiusmod high life accusamus terry
                    richardson ad squid. 3 wolf moon officia aute, non cupidatat skateboard dolor brunch.Food truck quinoa
                    nesciunt laborum eiusmod.Brunch 3 wolf moon tempor, sunt aliqua
                    put a bird on it squid single- origin coffee nulla assumenda
                    shoreditch et.Nihil anim keffiyeh helvetica, craft beer labore
                    wes anderson cred nesciunt sapiente ea proident.Ad vegan
                    excepteur butcher vice lomo.Leggings occaecat craft beer
                    farm - to - table, raw denim aesthetic synth nesciunt you probably
                    haven't heard of them accusamus labore sustainable VHS.`,
        },
    ];

    return (
        <Row>
            <Col lg={6}>
                <Card>
                    <Card.Body>
                        <h5 className="header-title mb-3 mt-0">Accordions Example</h5>

                        <Accordion defaultActiveKey="0" id="accordion" className="custom-accordionwitharrow">
                            {(accordianContent || []).map((item, index) => {
                                return (
                                    <CustomAccordion1
                                        key={index}
                                        item={item}
                                        index={index}
                                        length={accordianContent.length}
                                    />
                                );
                            })}
                        </Accordion>
                    </Card.Body>
                </Card>
            </Col>

            <Col lg={6}>
                <Card>
                    <Card.Body>
                        <h5 className="header-title mb-4 mt-0">Custom Accordion</h5>

                        <Accordion defaultActiveKey="0" id="accordion" className="custom-accordion">
                            {(accordianContent2 || []).map((item, index) => {
                                return (
                                    <CustomAccordion2
                                        key={index}
                                        item={item}
                                        index={index}
                                        length={accordianContent2.length}
                                    />
                                );
                            })}
                        </Accordion>
                    </Card.Body>
                </Card>
            </Col>
        </Row>
    );
};

export default Accordions;
