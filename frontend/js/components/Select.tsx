import React from 'react';

export type Option = {
  value: string
  label: string
}

type Props = {
  name: string
  value: string
  className: string
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void
  options: Option[]
}

export const Select: React.FC<Props> = ({name, value, className, onChange, options}) => {
    return (
      <select
        name={name}
        value={value}
        className={className}
        onChange={onChange}
        style={{width:"auto"}}
      >
        {options.map((o: Option) => {
          return (
            <option key={o.value} value={o.value}>{o.label}</option>
          )
        })}
      </select>
    );
}
// export default class Select extends React.Component {
//
//   constructor(props) {
//     super(props);
//   }
//
//   render() {
//     let options = [];
//     for (const o of this.props.options) {
//       options.push(<option key={o.value} value={o.value}>{o.label}</option>);
//     }
//
//     return (
//       <select
//         name={this.props.name}
//         value={this.props.value}
//         className={this.props.className}
//         onChange={this.props.onChange}
//         style={{width:"auto"}}
//       >
//         {options}
//       </select>
//     );
//   }
// }
//
// Select.propTypes = {
//   name: PropTypes.string.isRequired,
//   value: PropTypes.any.isRequired,
//   options: PropTypes.array.isRequired,
//   onChange: PropTypes.func.isRequired,
//   className: PropTypes.string,
// };
