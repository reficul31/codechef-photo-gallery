import 'core-js/fn/object/values' // shim Object.values for Chromium<54
import { createStore, applyMiddleware, compose, combineReducers } from 'redux'
import { createEpicMiddleware, combineEpics } from 'redux-observable'
import thunk from 'redux-thunk'

import reducer from './reducers'
import * as epics from './epics'

const rootReducer = combineReducers({
	album: reducer.album,
	photo: reducer.photo,
	user: reducer.user,
	showcase: reducer.showcase,
})

const rootEpic = combineEpics(
	...Object.values(epics)
	)

export default function configureStore({ReduxDevTools = undefined} = {}) {
	const enhancers = [
		applyMiddleware(
			createEpicMiddleware(rootEpic),
			thunk
			),
	]
	if (ReduxDevTools) {
		enhancers.push(ReduxDevTools.instrument())
	}
	const enhancer = compose(...enhancers)

	const store = createStore(
		rootReducer,
		undefined,
		enhancer
		)

	return store
}
