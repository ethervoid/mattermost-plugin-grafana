export function handleImageUpdate(store) {
    return (msg) => {
        if (!msg.data) {
            return;
        }

        store.dispatch({
            type: '_received_image_update',
            data: {...msg.data},
        });
    };
}