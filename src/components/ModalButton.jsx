import React, { Component, PropTypes } from 'react'

// import Modal from 'react-modal'

class ModalButton extends Component {

  static propTypes = {
    modal: PropTypes.func.isRequired,
    // modalProps: PropTypes.object,
    children: PropTypes.node.isRequired,
  }

  state = {
    showModal: false,
  }

  handleOnClick = () => (
    this.setState({ showModal: !this.state.showModal })
  )

  renderChildren = () => React.cloneElement(
    React.Children.only(this.props.children),
    { onClick: this.handleOnClick }
  )

  render() {
    // const ModalComponent = this.props.modal
    return (
      <span>
        {this.renderChildren()}
      </span>
    )
  }
}

/*
   <Modal
   className="ModalContainer"
   show={this.state.showModal}
   onHide={this.handleOnClick}
   backdropClassName="ModalContainer__backdrop"
   transition={transition}
   dialogTransitionTimeout={TRANSITION_TIMEOUT}
   backdropTransitionTimeout={TRANSITION_TIMEOUT}
   >
   <div className="ModalContainer__dialog">
   <ModalComponent onClose={this.handleOnClick} />
   </div>
   </Modal>
 */

/*
   <Modal
   isOpen={this.state.showModal}
   onRequestClose={this.handleOnClick}
   >
   <ModalComponent />
   </Modal>
 */

export default ModalButton