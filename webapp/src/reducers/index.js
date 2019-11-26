import {combineReducers} from 'redux';

const initialState = {
    image: null
}

function updatePanel(state = initialState, action) {
    switch (action.type) {
    case '_received_image_update':
        return {
          ...state,
          image: action.data.image
        };
    default:
        return state;
    }
}

export default combineReducers({updatePanel});