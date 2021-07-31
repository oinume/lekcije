import React, { useState } from 'react';
import { Option, Select } from '../Select';
import { sprintf } from 'sprintf-js';
import { range } from 'lodash-es';

// TODO: must be class. Add method isZero() and parse(). To be defined in another file.
export type NotificationTimeSpan = {
  fromHour: number;
  fromMinute: number;
  toHour: number;
  toMinute: number;
};

type Props = {
  editable: boolean;
  timeSpans: Array<any>;
  handleAdd: () => void;
  handleDelete: (index: number) => void;
  handleUpdate: () => void;
  handleOnChange: (name: string, index: number, value: any) => void;
  handleSetEditable: (editable: boolean) => void;
};

export const NotificationTimeSpanForm: React.FC<Props> = ({
  editable,
  timeSpans,
  handleAdd,
  handleDelete,
  handleUpdate,
  handleOnChange,
  handleSetEditable,
}) => {
  const onClickPlus = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>): void => {
    event.preventDefault();
    handleAdd();
  };

  const onClickMinus = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>, index: number): void => {
    event.preventDefault();
    handleDelete(index);
  };

  const onChange = (event: React.ChangeEvent<HTMLSelectElement>): void => {
    const [name, index] = event.target.name.split('_');
    handleOnChange(name, Number(index), event.target.value);
  };

  const createTimeSpanItem = (timeSpan: any, index: number) => {
    const hourOptions = range(0, 24).map((v) => {
      return { value: String(v), label: String(v) };
    });
    const minuteOptions = [0, 30].map((v) => {
      return { value: String(v), label: sprintf('%02d', v) };
    });

    return (
      // TODO: Use <ul><li>
      <p style={{ marginBottom: '0px' }}>
        <Select
          name={'fromHour_' + index}
          value={timeSpan.fromHour}
          onChange={onChange}
          options={hourOptions}
          className="custom-select custom-select-sm"
        />
        時 &nbsp;
        <Select
          name={'fromMinute_' + index}
          value={timeSpan.fromMinute}
          onChange={onChange}
          options={minuteOptions}
          className="custom-select custom-select-sm"
        />
        分 &nbsp; 〜 &nbsp;&nbsp;
        <Select
          name={'toHour_' + index}
          value={timeSpan.toHour}
          onChange={onChange}
          options={hourOptions}
          className="custom-select custom-select-sm"
        />
        時 &nbsp;
        <Select
          name={'toMinute_' + index}
          value={timeSpan.toMinute}
          onChange={onChange}
          options={minuteOptions}
          className="custom-select custom-select-sm"
        />
        分
        <button type="button" className="btn btn-link button-plus" onClick={(event) => onClickPlus(event)}>
          <i className="fas fa-plus-circle button-plus" aria-hidden="true" />
        </button>
        <button type="button" className="btn btn-link button-plus" onClick={(event) => onClickMinus(event, index)}>
          {/* TODO: no need to pass index */}
          <i className="fas fa-minus-circle button-plus" aria-hidden="true" />
        </button>
      </p>
    );
  };

  let content = [];
  if (editable) {
    if (timeSpans.length > 0) {
      timeSpans.map((timeSpan, i) => {
        content.push(createTimeSpanItem(timeSpan, i));
      });
    } else {
      handleAdd();
    }
  } else {
    if (timeSpans.length > 0) {
      for (let timeSpan of timeSpans) {
        content.push(
          <p>
            {timeSpan.fromHour}:{sprintf('%02d', timeSpan.fromMinute)} 〜 {timeSpan.toHour}:
            {sprintf('%02d', timeSpan.toMinute)}
          </p>
        );
      }
    } else {
      content.push(<p>データがありません。編集ボタンで追加できます。</p>);
    }
  }

  return (
    <form className="form-horizontal">
      <div className="form-group">
        <div className="col-sm-3">
          <label htmlFor="notificationTimeSpan" className="control-label">
            レッスン希望時間帯
          </label>
          <a href="https://lekcije.amebaownd.com/posts/3614832" target="_blank">
            <i className="fas fa-question-circle button-help" aria-hidden="true" />
          </a>
        </div>
        <div className="col-sm-7">{content}</div>
      </div>
      <div className="form-group">
        <div className="col-sm-offset-2 col-sm-8">
          <UpdateButton
            editable={editable}
            handleOnClick={(_) => {
              if (editable) {
                handleSetEditable(false);
                handleUpdate();
              } else {
                handleSetEditable(true);
              }
            }}
          />
        </div>
      </div>
    </form>
  );
};

type UpdateButtonProps = {
  editable: boolean;
  handleOnClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
};

const UpdateButton = ({ editable, handleOnClick }: UpdateButtonProps) => {
  const className = editable ? 'btn btn-primary' : 'btn btn-outline-primary';
  return (
    <button
      type="button"
      className={className}
      onClick={(event) => {
        event.preventDefault();
        handleOnClick(event);
      }}
    >
      {editable ? '更新' : '編集'}
    </button>
  );
};
