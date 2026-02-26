import { lazy, Suspense } from 'react';
import { Spin } from 'antd';

// 首页
const Home = lazy(() => import("../pages/home").then(m => ({ default: m.Home })));

// 登录
const Login = lazy(() => import("../pages/login").then(m => ({ default: m.Login })));
const Error = lazy(() => import("../utils/Error"));

// 告警管理
const AlertRuleList = lazy(() => import("../pages/alert/rule").then(m => ({ default: m.AlertRuleList })));
const AlertRule = lazy(() => import("../pages/alert/rule/create"));
const AlertRuleGroup = lazy(() => import("../pages/alert/ruleGroup").then(m => ({ default: m.AlertRuleGroup })));
const RuleTemplate = lazy(() => import("../pages/alert/tmpl").then(m => ({ default: m.RuleTemplate })));
const RuleTemplateGroup = lazy(() => import("../pages/alert/tmplGroup").then(m => ({ default: m.RuleTemplateGroup })));
const AlertSimulator = lazy(() => import("../pages/alert/simulator").then(m => ({ default: m.AlertSimulator })));
const Silences = lazy(() => import("../pages/silence").then(m => ({ default: m.Silences })));

// 通知管理
const NoticeObjects = lazy(() => import("../pages/notice").then(m => ({ default: m.NoticeObjects })));
const NoticeTemplate = lazy(() => import("../pages/notice/tmpl").then(m => ({ default: m.NoticeTemplate })));
const NoticeRecords = lazy(() => import("../pages/notice/history").then(m => ({ default: m.NoticeRecords })));

// 值班管理
const DutyManage = lazy(() => import("../pages/duty").then(m => ({ default: m.DutyManage })));
const CalendarApp = lazy(() => import("../pages/duty/calendar").then(m => ({ default: m.CalendarApp })));

// 数据源
const Datasources = lazy(() => import("../pages/datasources").then(m => ({ default: m.Datasources })));

// 仪表盘
const DashboardFolder = lazy(() => import("../pages/dashboards/folder").then(m => ({ default: m.DashboardFolder })));
const Dashboards = lazy(() => import("../pages/dashboards/dashboard").then(m => ({ default: m.Dashboards })));
const GrafanaDashboardComponent = lazy(() => import("../pages/dashboards/dashboard/iframe").then(m => ({ default: m.GrafanaDashboardComponent })));

// 用户管理
const User = lazy(() => import("../pages/members/user").then(m => ({ default: m.User })));
const UserRole = lazy(() => import("../pages/members/role").then(m => ({ default: m.UserRole })));

// 租户管理
const Tenants = lazy(() => import("../pages/tenant").then(m => ({ default: m.Tenants })));
const TenantDetail = lazy(() => import("../pages/tenant/detail").then(m => ({ default: m.TenantDetail })));

// 告警订阅
const Subscribe = lazy(() => import("../pages/subscribe").then(m => ({ default: m.Subscribe })));
const CreateSubscribeModel = lazy(() => import("../pages/subscribe/create").then(m => ({ default: m.CreateSubscribeModel })));

// 网络分析
const Probing = lazy(() => import("../pages/probing").then(m => ({ default: m.Probing })));
const CreateProbingRule = lazy(() => import("../pages/probing/create").then(m => ({ default: m.CreateProbingRule })));
const OnceProbing = lazy(() => import("../pages/probing/once").then(m => ({ default: m.OnceProbing })));

// 故障中心
const FaultCenter = lazy(() => import("../pages/faultCenter").then(m => ({ default: m.FaultCenter })));
const FaultCenterDetail = lazy(() => import("../pages/faultCenter/detail").then(m => ({ default: m.FaultCenterDetail })));

// 工单管理
const TicketList = lazy(() => import("../pages/ticket").then(m => ({ default: m.TicketList })));
const TicketDetail = lazy(() => import("../pages/ticket/detail").then(m => ({ default: m.TicketDetail })));
const TicketCreate = lazy(() => import("../pages/ticket/create").then(m => ({ default: m.TicketCreate })));
const RepairForm = lazy(() => import("../pages/ticket/repair").then(m => ({ default: m.default || m.RepairForm })));
const TicketReview = lazy(() => import("../pages/ticket/review").then(m => ({ default: m.default || m.TicketReview })));
const ReviewerManage = lazy(() => import("../pages/ticket/reviewer").then(m => ({ default: m.default || m.ReviewerManage })));
const TicketStatistics = lazy(() => import("../pages/ticket/statistics").then(m => ({ default: m.TicketStatistics })));
const TicketSlaPolicy = lazy(() => import("../pages/ticket/slaPolicy").then(m => ({ default: m.default || m.TicketSlaPolicy })));
const TicketTemplate = lazy(() => import("../pages/ticket/template").then(m => ({ default: m.default || m.TicketTemplate })));
const WorkHoursStandard = lazy(() => import("../pages/ticket/workHours").then(m => ({ default: m.default || m.WorkHoursStandard })));

// 知识库
const KnowledgeList = lazy(() => import("../pages/knowledge/index").then(m => ({ default: m.KnowledgeList })));
const KnowledgeDetail = lazy(() => import("../pages/knowledge/detail").then(m => ({ default: m.default || m.KnowledgeDetail })));
const KnowledgeCreate = lazy(() => import("../pages/knowledge/create").then(m => ({ default: m.default || m.KnowledgeCreate })));

// 智能派单
const AssignmentRule = lazy(() => import("../pages/assignment_rule/index").then(m => ({ default: m.AssignmentRule })));

// 系统设置
const AuditLog = lazy(() => import("../pages/audit").then(m => ({ default: m.AuditLog })));
const SystemSettings = lazy(() => import("../pages/settings").then(m => ({ default: m.SystemSettings })));
const Profile = lazy(() => import("../pages/profile").then(m => ({ default: m.default })));

// 移动端
const MobileTicketCreate = lazy(() => import("../pages/mobile/ticket/create").then(m => ({ default: m.MobileTicketCreate })));
const MobileTicketQuery = lazy(() => import("../pages/mobile/ticket/query").then(m => ({ default: m.MobileTicketQuery })));

// 组件
const ComponentsContent = lazy(() => import('../components').then(m => ({ default: m.ComponentsContent })));

// Loading组件
const PageLoading = () => (
	<div style={{ 
		display: 'flex', 
		justifyContent: 'center', 
		alignItems: 'center', 
		height: '200px' 
	}}>
		<Spin size="large" tip="加载中..." />
	</div>
);

// eslint-disable-next-line import/no-anonymous-default-export
export default [
    {
        path: '/',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="off" c={<Home />} /></Suspense>,
    },
    {
        path: '/login',
        element: <Suspense fallback={<PageLoading />}><Login /></Suspense>
    },
    {
        path: '/ruleGroup',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="告警规则组" c={<AlertRuleGroup />} /></Suspense>
    },
    {
        path: '/ruleGroup/:id/rule/list',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="告警规则" c={<AlertRuleList />} /></Suspense>
    },
    {
        path: '/ruleGroup/:id/rule/add',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="添加告警规则" c={<AlertRule type="add"/>} /></Suspense>
    },
    {
        path: '/ruleGroup/:id/rule/:ruleId/edit',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="编辑告警规则" c={<AlertRule type="edit"/>} /></Suspense>
    },
    {
        path: '/alert/simulator',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="告警模拟器" c={<AlertSimulator />} /></Suspense>
    },
    {
        path: '/silenceRules',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="静默规则" c={<Silences />} /></Suspense>
    },
    {
        path: '/tmplType/:tmplType/group',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="规则模版组" c={<RuleTemplateGroup />} /></Suspense>
    },
    {
        path: '/tmplType/:tmplType/:ruleGroupName/templates',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="规则模版" c={<RuleTemplate />} /></Suspense>
    },
    {
        path: '/noticeObjects',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="通知对象" c={<NoticeObjects />} /></Suspense>
    },
    {
        path: '/noticeTemplate',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="通知模版" c={<NoticeTemplate />} /></Suspense>
    },
    {
        path: '/noticeRecords',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="通知记录" c={<NoticeRecords />} /></Suspense>
    },
    {
        path: '/dutyManage',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="值班中心" c={<DutyManage />} /></Suspense>
    },
    {
        path: '/dutyManage/:id/calendar',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="值班表" c={<CalendarApp />} /></Suspense>
    },
    {
        path: '/user',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="用户管理" c={<User />} /></Suspense>
    },
    {
        path: '/userRole',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="角色管理" c={<UserRole />} /></Suspense>
    },
    {
        path: '/tenants',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="租户管理" c={<Tenants />} /></Suspense>
    },
    {
        path: '/tenants/detail/:id',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="租户" c={<TenantDetail/>} /></Suspense>
    },
    {
        path: '/datasource',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="数据源" c={<Datasources />} /></Suspense>
    },
    {
        path: '/folders',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="仪表盘目录" c={<DashboardFolder />} /></Suspense>
    },
    {
        path: '/folder/:id/list',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="仪表盘" c={<Dashboards />} /></Suspense>
    },
    {
        path: 'dashboard/f/:fid/g/:did/info',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="仪表盘详情" c={<GrafanaDashboardComponent />} /></Suspense>
    },
    {
        path: '/auditLog',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="日志审计" c={<AuditLog />} /></Suspense>
    },
    {
        path: '/settings',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="系统设置" c={<SystemSettings/>}/></Suspense>
    },
    {
        path: '/onceProbing',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="及时拨测" c={<OnceProbing/>} /></Suspense>
    },
    {
        path: '/probing',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="拨测任务" c={<Probing/>} /></Suspense>
    },
    {
        path: '/probing/create',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="创建拨测规则" c={<CreateProbingRule type="add"/>} /></Suspense>
    },
    {
        path: '/probing/:id/edit',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="编辑拨测规则" c={<CreateProbingRule type="edit"/>} /></Suspense>
    },
    {
        path: '/subscribes',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="告警订阅" c={<Subscribe />} /></Suspense>
    },
    {
        path: '/subscribe/create',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="添加订阅" c={<CreateSubscribeModel />} /></Suspense>
    },
    {
        path: '/profile',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="个人信息" c={<Profile />} /></Suspense>
    },
    {
        path: '/faultCenter',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="故障中心" c={<FaultCenter />} /></Suspense>
    },
    {
        path: '/faultCenter/detail/:id',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="故障中心详情" c={<FaultCenterDetail />} /></Suspense>
    },
    {
        path: '/ticket',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="工单管理" c={<TicketList />} /></Suspense>
    },
    {
        path: '/ticket/create',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="创建工单" c={<TicketCreate />} /></Suspense>
    },
    {
        path: '/ticket/repair',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="人工报修" c={<RepairForm />} /></Suspense>
    },
    {
        path: '/ticket/detail/:id',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="工单详情" c={<TicketDetail />} /></Suspense>
    },
{
        path: '/ticket/review',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="工单评审" c={<TicketReview />} /></Suspense>
    },
    {
        path: '/ticket/reviewer',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="评委管理" c={<ReviewerManage />} /></Suspense>
    },
    {
        path: '/ticket/statistics',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="工单统计" c={<TicketStatistics />} /></Suspense>
    },
    {
        path: '/ticket/sla',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="SLA策略" c={<TicketSlaPolicy />} /></Suspense>
    },
    {
        path: '/ticket/template',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="工单模板" c={<TicketTemplate />} /></Suspense>
    },
    {
        path: '/mobile/ticket/create',
        element: <Suspense fallback={<PageLoading />}><MobileTicketCreate /></Suspense>
    },
    {
        path: '/mobile/ticket/query',
        element: <Suspense fallback={<PageLoading />}><MobileTicketQuery /></Suspense>
    },
    {
        path: '/knowledge',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="知识库" c={<KnowledgeList />} /></Suspense>
    },
    {
        path: '/knowledge/detail/:id',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="知识详情" c={<KnowledgeDetail />} /></Suspense>
    },
    {
        path: '/knowledge/create',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="创建知识" c={<KnowledgeCreate />} /></Suspense>
    },
    {
        path: '/ticket/workHours',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="工时标准" c={<WorkHoursStandard />} /></Suspense>
    },
    {
        path: '/assignment-rule',
        element: <Suspense fallback={<PageLoading />}><ComponentsContent name="智能派单规则" c={<AssignmentRule />} /></Suspense>
    },
    {
        path: '/*',
        element: <Suspense fallback={<PageLoading />}><Error /></Suspense>
    }
]