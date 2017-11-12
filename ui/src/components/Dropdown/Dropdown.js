import React from 'react';
import { render, unmountComponentAtNode } from 'react-dom';
import * as dropdown from './heplpers';
import cx from 'classnames';

/**
 * It is a simple Dropdown.
 *
 * <Dropdown trigger={ <button>Trigger</button> }>
 *     Content of Dropdown.
 * </Dropdown>
 *
 */
export class Dropdown extends React.Component {


    static get defaultProps () {
        return {
            disabled: false,
            trigger: null,
            triggerClassName: null,
            contentClassName: null,
            isClosableOnTriggerClick: true,
            isClosableOnContentClick: false
        };
    }

    constructor (props) {
        super(props);
        this.state = {
            active: false
        };
        this._trigger = null;
        this._content = null;
        this._contentEl = this._createContentEl();
        this._triggerEl = null;
    }

    componentDidMount () {
        this._renderContent();
    }

    componentWillUnmount () {
        unmountComponentAtNode(this._contentEl);
        document.body.removeChild(this._contentEl);
    }

    componentDidUpdate () {
        this._renderContent();
    }

    isActive () {
        return (typeof this.props.active === 'boolean') ?
            this.props.active :
            this.state.active;
    }

    hide () {
        const {onHide} = this.props;
        this.setState({
            active: false
        });
        if (onHide) {
            onHide();
        }
    }

    show () {
        const {onShow} = this.props;
        this.setState({
            active: true
        });
        if (onShow) {
            onShow();
        }
    }

    _createContentEl () {
        const {isClosableOnContentClick} = this.props;
        const contentEl = document.createElement('div');
        contentEl.addEventListener('hide', this._handleDropdownHide.bind(this));
        if (isClosableOnContentClick) {
            contentEl.addEventListener('click', () => this.hide());
        }
        document.body.appendChild(contentEl);
        return contentEl;
    }

    _handleTriggerClick (event) {
        this._triggerEl = event.currentTarget;
        if (this.isActive()) {
            this.hide();
        } else {
            this.show();
        }
    }

    _handleDropdownHide (event) {
        if (this.isActive()) {
            this.hide();
        }
    }

    _getTrigger () {
        const {
            trigger,
            triggerClassName,
            disabled
        } = this.props;
        return React.cloneElement(this._trigger ? this._trigger : trigger, {
            className: cx(
                triggerClassName,
                this._trigger ? this._trigger.props.className: null,
                trigger ? trigger.props.className: null
            ),
            onClick: this._handleTriggerClick.bind(this),
            disabled: disabled
        });
    }

    _getContent () {
        if (this._content) {
            return React.cloneElement(this._content, {
                className: cx(this._content.props.className, this.props.contentClassName)
            });
        }
        return (
            <DropdownContent className={ this.props.contentClassName }>
                { this.props.children }
            </DropdownContent>
        );
    }

    _renderContent () {
        const {disabled} = this.props;
        if (disabled) {
            return;
        }
        const content = this._getContent();
        render(content, this._contentEl);
        this._contentEl.className = cx('dropdown-content', content.props.className);
        if (this.isActive()) {
            dropdown.show(this._triggerEl, this._contentEl);
        } else {
            dropdown.hide(this._triggerEl, this._contentEl);
        }
    }

    render () {
        React.Children.forEach(this.props.children, (child) => {
            switch (child.type) {
                case DropdownTrigger:
                    this._trigger = child;
                    break;
                case DropdownContent:
                    this._content = child;
                    break;
            }
        });
        return this._getTrigger();
    }
}

export const  DropdownContent =({children}) => (
    <div>
        { children }
    </div>
);


export class DropdownTrigger extends React.Component {

    _handleClick (event) {
        const {onClick} = this.props;
        event.preventDefault();
        if (onClick) {
            onClick(event);
        }
    }

    render () {
        const { children, className, disabled } = this.props;
        return React.cloneElement(children, {
            className: cx(children.props.className, 'dropdown-trigger', className, { disabled }),
            onClick: this._handleClick.bind(this),
            disabled: disabled
        });
    }
}
