package constant

// TextPointConstant 积分常量
const TextPointConstant = 3

const ImagePointConstant = 4

const VideoPointConstant = 5

/*
热度权重
*/

// LikeWeightConstant 点赞
const LikeWeightConstant = 0.5

// CollectWeightConstant 收藏
const CollectWeightConstant = 0.2

// CommentWeightConstant 评论
const CommentWeightConstant = 0.3

/*
互动类型
*/

// ArticleInteractionConstant 文章互动
const ArticleInteractionConstant = 0

// CommentInteractionConstant 评论互动
const CommentInteractionConstant = 1

/*
互动消息通知类型
*/

const LikeMsgConstant = 0

const CommentMsgConstant = 1

const CollectMsgConstant = 2

/*
广播消息通知类型
*/

const SystemMsgConstant = 1

const ManagerMsgConstant = 2

/*
文章字数限制
*/

const WordLimitMin = 0

const WordLimitMax = 300

/*
管理员推选人数限制
*/

const PeopleLimitClass = 3
const PeopleLimitGrade = 5
const PeopleLimitCollege = 10

/*
超级管理员用户id
*/

const SupserManId = 182

/*
一天对多发的话题数
*/

const ArticlePublishLimit = 10000

/*
优秀文章各级指数
*/

const ArticleClass = 1
const ArticleGrade = 3
const ArticleCollege = 5

/*
优秀帖子等级
*/

const CommonArticle = 0

const ClassArticle = 1

const GradeArticle = 2

const CollegeArticle = 3

/*
文章上传，文件大小限制
*/

const MemoryLimit = 83886080
