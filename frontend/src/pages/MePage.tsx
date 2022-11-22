import React, {useState} from 'react';
import {toast} from 'react-toastify';
import {useQueryClient} from '@tanstack/react-query';
import {PageTitle} from '../components/PageTitle';
import {Loader} from '../components/Loader';
import type {Teacher} from '../models/Teacher';
import {ToastContainer} from '../components/ToastContainer';
import type {GetViewerWithFollowingTeachersQuery} from '../graphql/generated';
import {
  useCreateFollowingTeacherMutation, useDeleteFollowingTeachersMutation,
  useGetViewerWithFollowingTeachersQuery, useGetViewerWithNotificationTimeSpansQuery,
} from '../graphql/generated';
import type {GraphQLError} from '../http/graphql';
import {createGraphQLClient, toMessage} from '../http/graphql';
import type {FollowingTeacher} from '../models/FollowingTeacher';
import {SubmitButton} from '../components/SubmitButton';

const graphqlClient = createGraphQLClient();

export const MePage = () => {
  const getViewerResult = useGetViewerWithFollowingTeachersQuery<GetViewerWithFollowingTeachersQuery, GraphQLError>(graphqlClient, {}, {
    onError(error) {
      toast.error(toMessage(error, 'データの取得に失敗しました'));
    },
  });
  const showTutorial = getViewerResult.data ? getViewerResult.data.viewer.showTutorial : false;
  const followingTeachers: FollowingTeacher[] = getViewerResult.data ? getViewerResult.data.viewer.followingTeachers.nodes.map(node => ({
    teacher: {
      id: node.teacher.id,
      name: node.teacher.name,
    },
  })) : [];

  return (
    <div id="followingForm">
      <ToastContainer
        closeOnClick={false}
      />
      <PageTitle>フォローしている講師</PageTitle>
      {
        getViewerResult.isLoading
          ? <Loader isLoading={getViewerResult.isLoading}/>
          : <MeContent followingTeachers={followingTeachers} showTutorial={showTutorial}/>
      }
    </div>
  );
};

type MeContentProps = {
  followingTeachers: FollowingTeacher[];
  showTutorial: boolean; // eslint-disable-line react/boolean-prop-naming
};

// Help URL
// https://lekcije.amebaownd.com/posts/{{ if .IsUserAgentPC }}2044879{{ end }}{{ if .IsUserAgentSP }}1577091{{ end }}{{ if .IsUserAgentTablet }}1577091{{ end }}

const MeContent = ({followingTeachers, showTutorial}: MeContentProps) => (
  <>
    {showTutorial ? <Tutorial/> : <div/>}
    <CreateForm/>
    <TeacherList followingTeachers={followingTeachers}/>
  </>
);

const Tutorial = () => (
  <div className="alert alert-success alert-dismissible" role="alert">
    <button type="button" className="btn-close" data-bs-dismiss="alert" aria-label="Close"/>
    <h4><i className="bi bi-info-square-fill"/> 講師をフォローするには</h4>
    <ol>
      <li><a href="https://eikaiwa.dmm.com/list/" className="alert-link" target="_blank" rel="noreferrer">DMM英会話</a>でお気に入りの講師のページにアクセスしましょう</li>
      <li>講師のURLをコピーしましょう(<a href="https://lekcije.amebaownd.com/posts/1577091" className="alert-link" target="_blank" rel="noreferrer">ヘルプ</a>)</li>
      <li>URLを下の入力欄にペーストしてフォローしましょう</li>
      <li>フォローすると、その講師の空きレッスンがあった時にメールでお知らせします</li>
    </ol>
  </div>
);

const CreateForm = () => {
  const [teacherIdOrUrl, setTeacherIdOrUrl] = useState('');
  const [submitDisabled, setSubmitDisabled] = useState(true);
  const [submitLoading, setSubmitLoading] = useState(false);

  const queryClient = useQueryClient();

  const createFollowingTeacherMutation = useCreateFollowingTeacherMutation<GraphQLError>(graphqlClient, {
    async onSuccess() {
      await queryClient.invalidateQueries(useGetViewerWithFollowingTeachersQuery.getKey());
      setTeacherIdOrUrl('');
      setSubmitDisabled(true);
      setSubmitLoading(false);
      toast.success('講師をフォローしました！');
    },
    onError(error) {
      // eslint-disable-next-line @typescript-eslint/no-base-to-string, @typescript-eslint/restrict-template-expressions
      console.error(`useCreateFollowingTeacherMutation.onError: err=${error}`);
      toast.error(toMessage(error, '講師のフォローに失敗しました'));
    },
  });

  return (
    <form
      onSubmit={event => {
        event.preventDefault();
        setSubmitDisabled(true);
        setSubmitLoading(true);
        createFollowingTeacherMutation.mutate({input: {teacherIdOrUrl}});
      }}
    >
      <p>
        講師のURLまたはIDを入力してフォローします<a href="https://lekcije.amebaownd.com/posts/2044879" rel="noreferrer" target="_blank"><i className="fas fa-question-circle button-help" aria-hidden="true"/></a><br/>
        <small><a href="https://eikaiwa.dmm.com/" rel="noreferrer" target="_blank">DMM英会話で講師を検索</a></small>
      </p>
      <div className="input-group mb-3">
        <input
          required
          autoFocus
          id="teacherIdsOrUrl"
          type="text"
          className="form-control"
          name="teacherIdsOrUrl"
          placeholder="https://eikaiwa.dmm.com/teacher/index/492/"
          value={teacherIdOrUrl}
          onChange={event => {
            event.preventDefault();
            setSubmitDisabled(event.currentTarget.value === '');
            setTeacherIdOrUrl(event.currentTarget.value);
          }}
        />
        <span className="px-2"/>
        <SubmitButton
          disabled={submitDisabled}
          loading={submitLoading}
        >
          送信
        </SubmitButton>
      </div>
    </form>
  );
};

type TeacherListProps = {
  followingTeachers: FollowingTeacher[];
};

const TeacherList = ({followingTeachers}: TeacherListProps) => {
  const [checkedIds, setCheckedIds] = useState<string[]>([]);
  const [deleteSubmitDisabled, setDeleteSubmitDisabled] = useState(true);
  const [deleteSubmitLoading, setDeleteSubmitLoading] = useState(false);

  const handleCheckboxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const targetId = event.target.value;
    if (event.target.checked) {
      setCheckedIds([...checkedIds, targetId]);
      setDeleteSubmitDisabled(false);
    } else {
      const restIds = checkedIds.filter(id => id !== targetId);
      setCheckedIds(restIds);
      setDeleteSubmitDisabled(restIds.length === 0);
    }
  };

  const queryClient = useQueryClient();
  const deleteFollowingTeacherMutation = useDeleteFollowingTeachersMutation<GraphQLError>(graphqlClient, {
    async onSuccess() {
      await queryClient.invalidateQueries(useGetViewerWithFollowingTeachersQuery.getKey());
      toast.success('講師のフォローを解除しました');
      setDeleteSubmitDisabled(true);
      setDeleteSubmitLoading(false);
    },
    onError(error) {
      // eslint-disable-next-line @typescript-eslint/no-base-to-string, @typescript-eslint/restrict-template-expressions
      console.error(`deleteFollowingTeacherMutation.onError: err=${error}`);
      // toast.error(`講師のフォロー解除に失敗しました: ${error.message}`);
      toast.error(toMessage(error, '講師のフォロー解除に失敗しました'));
    },
  },
  );

  return (
    <div id="followingTeachers">
      <form
        onSubmit={event => {
          event.preventDefault();
          setDeleteSubmitDisabled(true);
          setDeleteSubmitLoading(true);
          deleteFollowingTeacherMutation.mutate({input: {teacherIds: checkedIds}});
        }}
      >
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th scope="col" className="col-md-1">
                <SubmitButton
                  disabled={deleteSubmitDisabled}
                  loading={deleteSubmitLoading}
                >
                  削除
                </SubmitButton>
              </th>
              <th scope="col" className="col-md-11">
                講師
              </th>
            </tr>
          </thead>
          <tbody>
            {followingTeachers.map(ft => <TeacherRow key={ft.teacher.id} teacher={ft.teacher} handleOnChange={handleCheckboxChange}/>)}
          </tbody>
        </table>
      </form>
    </div>
  );
};

type TeacherRowProps = {
  teacher: Teacher;
  handleOnChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
};

const TeacherRow = ({teacher, handleOnChange}: TeacherRowProps) => (
  <tr>
    <td className="col-md-1"><input type="checkbox" name="teacherIds" value={teacher.id} onChange={handleOnChange}/></td>
    <td className="col-md-8"><a href={`https://eikaiwa.dmm.com/teacher/index/${teacher.id}`} target="_blank" rel="noreferrer">{ teacher.name }</a></td>
  </tr>
);
