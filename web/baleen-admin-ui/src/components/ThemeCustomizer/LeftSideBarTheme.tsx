import React from 'react';
import { Form } from 'react-bootstrap';
import * as layoutConstants from '../../constants/layout';

interface LeftSideBarThemeProps {
    changeLeftSidebarTheme: (value: any) => void;
    leftSideBarTheme: string;
    layoutConstants: typeof layoutConstants.SideBarTheme;
}

const LeftSideBarTheme = ({ changeLeftSidebarTheme, leftSideBarTheme, layoutConstants }: LeftSideBarThemeProps) => {
    return (
        <>
            <h6 className="fw-medium mt-4 mb-2 pb-1">Left Sidebar Color</h6>

            <Form.Check className="form-check form-switch mb-1">
                <Form.Check.Input
                    type="radio"
                    name="theme"
                    id="light-check"
                    value={layoutConstants.LEFT_SIDEBAR_THEME_LIGHT}
                    onChange={(e) => changeLeftSidebarTheme(e.target.value)}
                    checked={leftSideBarTheme === layoutConstants.LEFT_SIDEBAR_THEME_LIGHT}
                />
                <Form.Check.Label htmlFor="light-check">Light</Form.Check.Label>
            </Form.Check>

            <Form.Check className="form-check form-switch mb-1">
                <Form.Check.Input
                    type="radio"
                    name="theme"
                    id="dark-check"
                    value={layoutConstants.LEFT_SIDEBAR_THEME_DARK}
                    onChange={(e) => changeLeftSidebarTheme(e.target.value)}
                    checked={leftSideBarTheme === layoutConstants.LEFT_SIDEBAR_THEME_DARK}
                />
                <Form.Check.Label htmlFor="dark-check">Dark</Form.Check.Label>
            </Form.Check>
        </>
    );
};

export default LeftSideBarTheme;
