import {id as pluginId} from './manifest';
import ChannelSubHeader from './components/channel_subheader';
import {handlePanelUpdate, handlePanelDeletion} from './websocket';
import Reducer from './reducers';

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerReducer(Reducer);
        registry.registerChannelSubHeaderComponent(ChannelSubHeader);
        registry.registerWebSocketEventHandler('custom_com.mattermost.grafana_update_grafana_subscription', handlePanelUpdate(store));
        registry.registerWebSocketEventHandler('custom_com.mattermost.grafana_remove_grafana_subscription', handlePanelDeletion(store));
    }
}

window.registerPlugin(pluginId, new Plugin());
