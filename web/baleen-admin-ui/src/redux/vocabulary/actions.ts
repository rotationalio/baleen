import { VocabulariesActionTypes } from './constants';

export const getAllVocabularies = () => ({
    type: VocabulariesActionTypes.GET_VOCABULARY,
    payload: {},
});

export const getAllVocabulariesResponseSuccess = (data: any) => ({
    type: VocabulariesActionTypes.GET_VOCABULARY_RESPONSE_SUCCESS,
    payload: { data },
});

export const getAllVocabulariesResponseError = (error: any) => ({
    type: VocabulariesActionTypes.GET_VOCABULARY_RESPONSE_SUCCESS,
    payload: { error },
});
