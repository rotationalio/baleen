import React from 'react'

export default function useSafeDispatch(dispatch) {
    const mounted = React.useRef(false)

    React.useLayoutEffect(() => {
        mounted.current = true
        return () => (mounted.current = false)
    }, [])

    return React.useCallback(
        (...args) => (mounted.current ? dispatch(...args) : void 0),
        [dispatch],
    )
}