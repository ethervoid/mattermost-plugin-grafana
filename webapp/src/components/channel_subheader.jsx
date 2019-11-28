import React from 'react';
import PropTypes from 'prop-types';

import {connect} from 'react-redux';

import {id as pluginId} from 'manifest';

class ChannelSubHeader extends React.PureComponent {
    static propTypes = {
        image: PropTypes.string.isRequired,
        channelId: PropTypes.string,
        channelTarget: PropTypes.string,
        show: PropTypes.bool.isRequired,
    };
    render() {
        if (this.props.image === null) {
            return null;
        }
        if (this.props.channelId !== this.props.channelTarget) {
            return null;
        }
        if (!this.props.show) {
            return null;
        }

        const image = 'data:image/png;base64,' + this.props.image;
        const style = {};

        return (
            <div
                className='grafana'
                style={style}
            >
                <img
                    key='grafanaPanel'
                    src={image}
                />
            </div>
        );
    }
}

function mapStateToProps(state) {
    return {
        image: state[`plugins-${pluginId}`].updatePanel.image,
        channelTarget: state[`plugins-${pluginId}`].updatePanel.channelTarget,
        show: state[`plugins-${pluginId}`].updatePanel.show,
    };
}

export default connect(mapStateToProps)(ChannelSubHeader);