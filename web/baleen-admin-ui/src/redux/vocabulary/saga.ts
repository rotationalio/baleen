import { VocabulariesActionTypes } from './constants';
import { AxiosResponse } from 'axios';
import { getAllVocabulariesResponseSuccess, getAllVocabulariesResponseError } from './actions';
import { call, put, takeEvery } from 'redux-saga/effects';
import { getVocabularies } from 'services/vocabulary';

function* getAllVocabularies() {
    try {
        const response: AxiosResponse = yield call(getVocabularies);
        yield put(getAllVocabulariesResponseSuccess(response.data));
    } catch (error) {
        yield put(getAllVocabulariesResponseError(error));
    }
}

export function* vocabulariesSaga() {
    yield takeEvery(VocabulariesActionTypes.GET_VOCABULARY, getAllVocabularies);
}
