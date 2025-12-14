"""
NLP Service 数据模型
"""
from pydantic import BaseModel, Field
from typing import Optional, List, Dict, Any
from datetime import datetime


# ============================================
# 文本提取相关
# ============================================
class TextExtractionRequest(BaseModel):
    """文本提取请求"""
    content: str = Field(..., description="HTML或纯文本内容")
    content_type: str = Field(default="html", description="内容类型: html/text/pdf")


class TextExtractionResponse(BaseModel):
    """文本提取响应"""
    extracted_text: str = Field(..., description="提取的文本")
    title: Optional[str] = Field(None, description="标题")
    metadata: Dict[str, Any] = Field(default_factory=dict, description="元数据")


# ============================================
# 实体识别相关
# ============================================
class EntityRecognitionRequest(BaseModel):
    """实体识别请求"""
    text: str = Field(..., description="待识别文本")
    entity_types: Optional[List[str]] = Field(
        default=["time", "location", "organization", "requirement"],
        description="需要识别的实体类型"
    )


class Entity(BaseModel):
    """实体"""
    text: str = Field(..., description="实体文本")
    type: str = Field(..., description="实体类型")
    start: int = Field(..., description="起始位置")
    end: int = Field(..., description="结束位置")
    confidence: float = Field(..., description="置信度")


class EntityRecognitionResponse(BaseModel):
    """实体识别响应"""
    entities: List[Entity] = Field(..., description="识别到的实体列表")
    structured_data: Dict[str, Any] = Field(default_factory=dict, description="结构化数据")


# ============================================
# OCR相关
# ============================================
class OCRRequest(BaseModel):
    """OCR请求"""
    image_url: Optional[str] = Field(None, description="图片URL")
    image_base64: Optional[str] = Field(None, description="Base64编码的图片")


class OCRResponse(BaseModel):
    """OCR响应"""
    text: str = Field(..., description="识别的文本")
    confidence: float = Field(..., description="平均置信度")
    regions: List[Dict[str, Any]] = Field(default_factory=list, description="文本区域")


# ============================================
# 向量化相关
# ============================================
class VectorizeRequest(BaseModel):
    """向量化请求"""
    text: str = Field(..., description="待向量化文本")
    model: str = Field(default="paraphrase-multilingual-MiniLM-L12-v2", description="使用的模型")


class VectorizeResponse(BaseModel):
    """向量化响应"""
    vector: List[float] = Field(..., description="文本向量")
    dimension: int = Field(..., description="向量维度")


# ============================================
# 机会结构化相关
# ============================================
class OpportunityStructureRequest(BaseModel):
    """机会结构化请求"""
    raw_html: str = Field(..., description="原始HTML内容")
    source_url: str = Field(..., description="来源URL")


class StructuredOpportunity(BaseModel):
    """结构化机会"""
    title: str
    type: str  # 竞赛/实习/项目/奖学金
    description: str
    organizer: Optional[str] = None
    start_date: Optional[str] = None
    deadline: Optional[str] = None
    event_date: Optional[str] = None
    location: Optional[str] = None
    requirements: Dict[str, Any] = Field(default_factory=dict)
    tags: List[str] = Field(default_factory=list)
    target_majors: List[str] = Field(default_factory=list)
    confidence: float = Field(default=0.0, description="结构化置信度")


class OpportunityStructureResponse(BaseModel):
    """机会结构化响应"""
    opportunity: StructuredOpportunity
    raw_entities: List[Entity] = Field(default_factory=list)
    processing_time: float = Field(..., description="处理时间(秒)")


# ============================================
# 竞赛级别识别相关
# ============================================
class CompetitionLevelRequest(BaseModel):
    """竞赛级别识别请求"""
    title: str = Field(..., description="竞赛标题")
    description: str = Field(..., description="竞赛描述")
    organizer: Optional[str] = Field(None, description="主办方")


class CompetitionLevelResponse(BaseModel):
    """竞赛级别识别响应"""
    level: str = Field(..., description="级别: 国家级A类/国家级B类/省级/校级/国际级/其他")
    confidence: float = Field(..., description="置信度")
    reason: str = Field(..., description="判定理由")
    certification_source: Optional[str] = Field(None, description="认定来源")
    organizer_type: Optional[str] = Field(None, description="主办方类型")


# ============================================
# 通用响应
# ============================================
class HealthResponse(BaseModel):
    """健康检查响应"""
    status: str
    service: str
    version: str
    timestamp: str


class ErrorResponse(BaseModel):
    """错误响应"""
    error: str
    detail: Optional[str] = None
    timestamp: str = Field(default_factory=lambda: datetime.now().isoformat())
