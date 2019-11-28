export function handlePanelUpdate(store) {
    return (msg) => {
        if (!msg.data) {
            return;
        }

        store.dispatch({
            type: '_received_panel_update',
            data: {...msg.data},
        });
    };
}

export function handlePanelDeletion(store) {
    return (msg) => {
        if (!msg.data) {
            return;
        }

        store.dispatch({
            type: '_received_panel_deletion',
            data: {...msg.data},
        });
    };
}