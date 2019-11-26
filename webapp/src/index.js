import {id as pluginId} from './manifest';
import ChannelSubHeader from './components/channel_subheader';
import {handleImageUpdate} from './websocket';
import Reducer from './reducers';

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerReducer(Reducer);
        registry.registerChannelSubHeaderComponent(ChannelSubHeader);
        registry.registerWebSocketEventHandler('custom_com.mattermost.grafana_update_grafana_subscription', handleImageUpdate(store));
    }
}

window.registerPlugin(pluginId, new Plugin());
