import React, {useState} from 'react';
import {sprintf} from 'sprintf-js';
import {range} from 'lodash-es';
import {Select} from '../Select';

// TODO: must be class. Add method isZero() and parse(). To be defined in another file.
export type NotificationTimeSpan = {
  fromHour: number;
  fromMinute: number;
  toHour: number;
  toMinute: number;
};

type Props = {
  timeSpans: NotificationTimeSpan[];
  handleAdd: () => void;
  handleDelete: (index: number) => void;
  handleUpdate: () => void;
  handleOnChange: (name: string, index: number, value: any) => void;
};

export const NotificationTimeSpanForm: React.FC<Props> = ({
  timeSpans,
  handleAdd,
  handleDelete,
  handleUpdate,
  handleOnChange,
}) => {
  const [editable, setEditable] = useState<boolean>(false);

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

  let content: any = [];
  if (timeSpans.length > 0) {
    content = timeSpans.map((timeSpan, i) => (
      <TimeSpanItem
        key={i}
        editable={editable}
        timeSpan={timeSpan}
        index={i}
        handleOnChange={onChange}
        handleOnClickPlus={onClickPlus}
        handleOnClickMinus={onClickMinus}
      />
    ));
  } else if (editable) {
    handleAdd();
  } else {
    content = [<p key="noTimeSpans">データがありません。編集ボタンで追加できます。</p>];
  }

  return (
    <form>
      <h5>
        レッスン希望時間帯
        {' '}
        <a href="https://lekcije.amebaownd.com/posts/3614832" target="_blank" rel="noreferrer">
          <i className="fas fa-question-circle button-help" aria-hidden="true"/>
        </a>
      </h5>
      <div className="form-group">
        <div className="col-sm-7">{content}</div>
      </div>
      <div className="form-group">
        <div className="col-sm-offset-2 col-sm-8">
          <UpdateButton
            editable={editable}
            handleOnClick={_ => {
              if (editable) {
                setEditable(false);
                handleUpdate();
              } else {
                setEditable(true);
              }
            }}
          />
        </div>
      </div>
    </form>
  );
};

type TimeSpanItemProps = {
  editable: boolean;
  timeSpan: any;
  index: number;
  handleOnChange: (event: React.ChangeEvent<HTMLSelectElement>) => void;
  handleOnClickPlus: (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => void;
  handleOnClickMinus: (event: React.MouseEvent<HTMLButtonElement, MouseEvent>, index: number) => void;
};

const TimeSpanItem = ({
  editable,
  timeSpan,
  index,
  handleOnChange,
  handleOnClickPlus,
  handleOnClickMinus,
}: TimeSpanItemProps) => {
  const hourOptions = range(0, 24).map(v => ({value: String(v), label: String(v)}));
  const minuteOptions = [0, 30].map(v => ({value: String(v), label: sprintf('%02d', v)}));

  if (editable) {
    return (
      // TODO: Use <ul><li>
      <p style={{marginBottom: '0px'}}>
        <Select
          name={`fromHour_${index}`}
          value={timeSpan.fromHour}
          options={hourOptions}
          className="custom-select custom-select-sm"
          onChange={handleOnChange}
        />
        時 &nbsp;
        <Select
          name={`fromMinute_${index}`}
          value={timeSpan.fromMinute}
          options={minuteOptions}
          className="custom-select custom-select-sm"
          onChange={handleOnChange}
        />
        分 &nbsp; 〜 &nbsp;&nbsp;
        <Select
          name={`toHour_${index}`}
          value={timeSpan.toHour}
          options={hourOptions}
          className="custom-select custom-select-sm"
          onChange={handleOnChange}
        />
        時 &nbsp;
        <Select
          name={`toMinute_${index}`}
          value={timeSpan.toMinute}
          options={minuteOptions}
          className="custom-select custom-select-sm"
          onChange={handleOnChange}
        />
        分
        <button
          type="button" className="btn btn-link button-plus" onClick={event => {
            handleOnClickPlus(event);
          }}
        >
          <i className="fas fa-plus-circle button-plus" aria-hidden="true"/>
        </button>
        <button
          type="button"
          className="btn btn-link button-plus"
          onClick={event => {
            handleOnClickMinus(event, index);
          }}
        >
          {/* TODO: no need to pass index */}
          <i className="fas fa-minus-circle button-plus" aria-hidden="true"/>
        </button>
      </p>
    );
  }

  return (
    <p>
      {timeSpan.fromHour}:{sprintf('%02d', timeSpan.fromMinute)} 〜 {timeSpan.toHour}:
      {sprintf('%02d', timeSpan.toMinute)}
    </p>
  );
};

type UpdateButtonProps = {
  editable: boolean;
  handleOnClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
};

const UpdateButton = ({editable, handleOnClick}: UpdateButtonProps) => {
  const className = editable ? 'btn btn-primary' : 'btn btn-outline-primary';
  return (
    <button
      type="button"
      className={className}
      onClick={event => {
        event.preventDefault();
        handleOnClick(event);
      }}
    >
      {editable ? '更新' : '編集'}
    </button>
  );
};
