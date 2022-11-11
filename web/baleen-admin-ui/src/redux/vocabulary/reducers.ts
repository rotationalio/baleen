import { VocabulariesActionTypes } from './constants';

type ActionType = {
    type: string;
    payload: Record<string, any>;
};

type State = {
    data: Record<string, any> | Array<any> | null;
    isLoading: boolean;
};
const INITIAL_STATE = {
    data: null,
    isLoading: false,
};

const vocabularies = (state: State = INITIAL_STATE, action: ActionType) => {
    switch (action.type) {
        case VocabulariesActionTypes.GET_VOCABULARY:
            return {
                ...state,
                isLoading: true,
            };
        case VocabulariesActionTypes.GET_VOCABULARY_RESPONSE_SUCCESS:
            return {
                ...state,
                isLoading: false,
                data: action.payload.data,
            };
        case VocabulariesActionTypes.GET_VOCABULARY_RESPONSE_ERROR:
            return {
                ...state,
                isLoading: false,
                data: action.payload.error,
            };

        default:
            return state;
    }
};

export { vocabularies };
