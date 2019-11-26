import React from 'react';
import {id as pluginId} from 'manifest';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';

class ChannelSubHeader extends React.PureComponent {
    static propTypes = {
        image: PropTypes.string.isRequired,
    };
    
    render() {
        const image = 'data:image/png;base64,' + this.props.image;
        const style = {}

        return (
            <div class="grafana" style={style}>
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
        image: state[`plugins-${pluginId}`].updatePanel.image
    };
}


export default connect(mapStateToProps)(ChannelSubHeader);