import React from "react";
import Select from '../Select';
import {sprintf} from 'sprintf-js';
import PropTypes from "prop-types";

export default class NotificationTimeSpanForm extends React.Component {
  constructor(props) {
    super(props);
    this.onClickPlus = this.onClickPlus.bind(this);
    this.onChange = this.onChange.bind(this);
  }

  onClickPlus(e) {
    e.preventDefault();
    this.props.dispatch('addTimeSpan');
  }

  onClickMinus(e, index) {
    e.preventDefault();
    this.props.dispatch('deleteTimeSpan', index);
  }

  onChange(event) {
    const a = event.target.name.split('_');
    const name = a[0];
    const index = a[1];
    const timeSpans = this.props.timeSpans.slice();
    timeSpans[index][name] = event.target.value;
    this.props.dispatch('onChangeTimeSpan', name, index, event.target.value);
  }

  createTimeSpanRow(timeSpan, index) {
    let hourOptions = [];
    for (let i = 0; i <= 23; i++) {
      hourOptions.push({value: i, label: i})
    }
    let minuteOptions = [];
    for (const i of [0, 30]) {
      minuteOptions.push({value: i, label: sprintf('%02d', i)});
    }

    return (
      // TODO: Use <ul><li>
      <p style={{marginBottom:"0px"}}>
        <Select
          name={'fromHour_' + index}
          value={timeSpan.fromHour}
          onChange={this.onChange}
          options={hourOptions}
          className="custom-select custom-select-sm"
        />時
        &nbsp;
        <Select
          name={'fromMinute_' + index}
          value={timeSpan.fromMinute}
          onChange={this.onChange}
          options={minuteOptions}
          className="custom-select custom-select-sm"
        />分

        &nbsp; 〜 &nbsp;&nbsp;

        <Select
          name={'toHour_' + index}
          value={timeSpan.toHour}
          onChange={this.onChange}
          options={hourOptions}
          className="custom-select custom-select-sm"
        />時
        &nbsp;
        <Select
          name={'toMinute_' + index}
          value={timeSpan.toMinute}
          onChange={this.onChange}
          options={minuteOptions}
          className="custom-select custom-select-sm"
        />分

        <button type="button" className="btn btn-link button-plus" onClick={() => this.onClickPlus(event)}>
          <i className="fas fa-plus-circle button-plus" aria-hidden="true" />
        </button>
        <button type="button" className="btn btn-link button-plus" onClick={() => this.onClickMinus(event, index)}>
          <i className="fas fa-minus-circle button-plus" aria-hidden="true" />
        </button>
      </p>
    );
  }

  render() {
    let content = [];
    if (this.props.editable) {
      if (this.props.timeSpans.length > 0) {
        this.props.timeSpans.map((timeSpan, i) => {
          content.push(this.createTimeSpanRow(timeSpan, i));
        });
      } else {
        this.props.dispatch('addTimeSpan');
      }
    } else {
      if (this.props.timeSpans.length > 0) {
        for (let timeSpan of this.props.timeSpans) {
          content.push(
            <p>{timeSpan.fromHour}:{sprintf('%02d', timeSpan.fromMinute)} 〜 {timeSpan.toHour}:{sprintf('%02d', timeSpan.toMinute)}</p>);
        }
      } else {
        content.push(
          <p>データがありません。編集ボタンで追加できます。</p>
        );
      }
    }

    let updateButton;
    if (this.props.editable) {
      updateButton =
        <button
          type="button" className="btn btn-primary"
          onClick={() => this.props.dispatch('updateTimeSpan', false)}
        >
          更新
        </button>;
    } else {
      updateButton =
        <button
          type="button" className="btn btn-outline-primary"
          onClick={() => this.props.dispatch('setTimeSpanEditable', true)}
        >
          編集
        </button>;
    }

    return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-3">
            <label htmlFor="notificationTimeSpan" className="control-label">レッスン希望時間帯</label>
            <a href="https://lekcije.amebaownd.com/posts/3614832" target="_blank">
              <i className="fas fa-question-circle button-help" aria-hidden="true" />
            </a>
          </div>
          <div className="col-sm-7">
            {content}
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
            {updateButton}
          </div>
        </div>
      </form>
    );
  }
}

NotificationTimeSpanForm.propTypes = {
  editable: PropTypes.bool.isRequired,
  timeSpans: PropTypes.array.isRequired,
  dispatch: PropTypes.func.isRequired,
};
