import React, { useCallback, useEffect, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';

// actions
import { hideRightSidebar } from '../redux/actions';

// store
import { AppDispatch, RootState } from '../redux/store';

// components
import ThemeCustomizer from '../components/ThemeCustomizer/';
import Scrollbar from '../components/Scrollbar';

interface RightSideBarProps {
    hideRightSidebar?: () => void;
    title?: string;
    children?: any;
}

const RightSideBar = (props: RightSideBarProps) => {
    const dispatch = useDispatch<AppDispatch>();
    const rightBarNodeRef: any = useRef(null);

    const { isOpenRightSideBar } = useSelector((state: RootState) => ({
        isOpenRightSideBar: state.Layout.isOpenRightSideBar,
    }));

    /**
     * Handle the click anywhere in doc
     */
    const handleOtherClick = useCallback(
        (e: any) => {
            if (isOpenRightSideBar) {
                if (rightBarNodeRef && rightBarNodeRef.current && rightBarNodeRef.current.contains(e.target)) {
                    return;
                } else {
                    dispatch(hideRightSidebar());
                }
            }
        },
        [rightBarNodeRef, dispatch, isOpenRightSideBar]
    );

    useEffect(() => {
        document.addEventListener('mousedown', handleOtherClick, false);
        return () => {
            document.removeEventListener('mousedown', handleOtherClick, false);
        };
    }, [handleOtherClick]);

    return (
        <React.Fragment>
            <div className="right-bar" ref={rightBarNodeRef}>
                <Scrollbar style={{ maxHeight: '100%', zIndex: 10000 }} timeout={500} scrollbarMaxSize={320}>
                    <ThemeCustomizer />
                </Scrollbar>
            </div>
            <div className="rightbar-overlay"></div>
        </React.Fragment>
    );
};

export default RightSideBar;
