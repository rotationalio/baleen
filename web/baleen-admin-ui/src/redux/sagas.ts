import { all } from 'redux-saga/effects';
import layoutSaga from './layout/saga';
import { vocabulariesSaga } from './vocabulary';

export default function* rootSaga() {
    yield all([layoutSaga(), vocabulariesSaga()]);
}
