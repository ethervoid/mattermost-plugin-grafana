import {combineReducers} from 'redux';

const initialState = {
    image: null,
    show: true,
    channelId: null,
};

function updatePanel(state = initialState, action) {
    switch (action.type) {
    case '_received_panel_update':
        return {
            ...state,
            image: action.data.image,
            channelTarget: action.data.channelTarget,
        };
    default:
        return state;
    }
}

function removePanel(state = initialState, action) {
    switch (action.type) {
    case '_received_panel_deletion':
        return {
            ...state,
            show: false,
            image: null,
            channelTarget: action.data.channelTarget,
        };
    default:
        return state;
    }
}

export default combineReducers({updatePanel, removePanel});